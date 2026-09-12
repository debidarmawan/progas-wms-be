package usecase

import (
	"fmt"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"time"
)

type DeliveryOrderUsecase interface {
	Issue(actorUserId string, req *dto.IssueDeliveryOrderRequest) (*dto.DeliveryOrderResponse, global.ErrorResponse)
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.DeliveryOrderResponse], global.ErrorResponse)
	FindById(id string) (*dto.DeliveryOrderResponse, global.ErrorResponse)
}

type deliveryOrderUsecase struct {
	txManager         helper.TxManager
	deliveryOrderRepo repository.DeliveryOrderRepository
	cylinderRepo      repository.CylinderRepository
	ledgerRepo        repository.CylinderLedgerRepository
	customerRepo      repository.CustomerRepository
	fleetRepo         repository.FleetRepository
	auditLogRepo      repository.AuditLogRepository
	pricingUsecase    CustomerItemPriceUsecase
	salesOrderRepo    repository.SalesOrderRepository
}

func NewDeliveryOrderUsecase(
	txManager helper.TxManager,
	deliveryOrderRepo repository.DeliveryOrderRepository,
	cylinderRepo repository.CylinderRepository,
	ledgerRepo repository.CylinderLedgerRepository,
	customerRepo repository.CustomerRepository,
	fleetRepo repository.FleetRepository,
	auditLogRepo repository.AuditLogRepository,
	pricingUsecase CustomerItemPriceUsecase,
	salesOrderRepo repository.SalesOrderRepository,
) DeliveryOrderUsecase {
	return &deliveryOrderUsecase{
		txManager:         txManager,
		deliveryOrderRepo: deliveryOrderRepo,
		cylinderRepo:      cylinderRepo,
		ledgerRepo:        ledgerRepo,
		customerRepo:      customerRepo,
		fleetRepo:         fleetRepo,
		auditLogRepo:      auditLogRepo,
		pricingUsecase:    pricingUsecase,
		salesOrderRepo:    salesOrderRepo,
	}
}

func (u *deliveryOrderUsecase) Issue(actorUserId string, req *dto.IssueDeliveryOrderRequest) (*dto.DeliveryOrderResponse, global.ErrorResponse) {
	if err := helper.ValidateBarcodeList(req.Barcodes); err != nil {
		return nil, global.BadRequestError(err.Error())
	}

	customer, err := u.customerRepo.FindById(req.CustomerId)
	if err != nil {
		return nil, err
	}
	if !customer.IsActive {
		return nil, global.BadRequestError("customer is not active")
	}

	fleet, err := u.fleetRepo.FindById(req.FleetId)
	if err != nil {
		return nil, err
	}
	if !fleet.IsActive {
		return nil, global.BadRequestError("fleet vehicle is not active")
	}

	var salesOrder *model.SalesOrder
	if req.SalesOrderId != "" {
		var salesOrderErr global.ErrorResponse
		salesOrder, salesOrderErr = u.salesOrderRepo.FindById(req.SalesOrderId)
		if salesOrderErr != nil {
			return nil, salesOrderErr
		}
		if salesOrder.CustomerId != customer.Id {
			return nil, global.BadRequestError("sales order customer does not match delivery order customer")
		}
		if salesOrder.Status != enum.SalesOrderStatusConfirmed && salesOrder.Status != enum.SalesOrderStatusPartial {
			return nil, global.BadRequestError("sales order must be confirmed before issuing a delivery order")
		}
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	cylinders, err := u.cylinderRepo.FindByBarcodes(tx, req.Barcodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if validationErr := helper.ValidateDOCylinders(cylinders); validationErr != nil {
		tx.Rollback()
		return nil, global.BadRequestError(validationErr.Error())
	}

	if salesOrder != nil {
		remainingByItem := make(map[string]int, len(salesOrder.Lines))
		lineByItem := make(map[string]*model.SalesOrderLine, len(salesOrder.Lines))
		for i := range salesOrder.Lines {
			remainingByItem[salesOrder.Lines[i].MasterItemId] = salesOrder.Lines[i].QtyOrdered - salesOrder.Lines[i].QtyDelivered
			lineByItem[salesOrder.Lines[i].MasterItemId] = &salesOrder.Lines[i]
		}
		for _, cylinder := range cylinders {
			if remainingByItem[cylinder.ItemId] <= 0 {
				tx.Rollback()
				return nil, global.BadRequestError("delivery quantity exceeds sales order remaining quantity")
			}
			remainingByItem[cylinder.ItemId]--
			lineByItem[cylinder.ItemId].QtyDelivered++
		}
	}

	totalWeight := helper.SumCylinderWeight(cylinders)
	if totalWeight > fleet.MaxWeightKg {
		tx.Rollback()
		return nil, global.BadRequestError(fmt.Sprintf(
			"total weight %.2f kg exceeds fleet max capacity %.2f kg",
			totalWeight, fleet.MaxWeightKg,
		))
	}

	order := &model.DeliveryOrder{
		BaseModel:     model.BaseModel{CreatedBy: actorUserId},
		DONumber:      helper.GenerateDONumber(),
		SalesOrderId:  req.SalesOrderId,
		CustomerId:    customer.Id,
		FleetId:       fleet.Id,
		Status:        enum.DeliveryOrderStatusInTransit,
		TotalWeightKg: totalWeight,
		CylinderQty:   len(cylinders),
		Notes:         req.Notes,
	}

	if err := u.deliveryOrderRepo.Create(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	details := make([]model.DeliveryOrderDetail, 0, len(cylinders))
	cylinderIds := make([]string, 0, len(cylinders))
	for _, cyl := range cylinders {
		weight := helper.CylinderFilledWeightKg(cyl.MasterItem)
		unitPrice, priceSource, priceErr := u.pricingUsecase.ResolvePrice(customer.Id, cyl.ItemId, time.Now())
		if priceErr != nil {
			tx.Rollback()
			return nil, priceErr
		}
		details = append(details, model.DeliveryOrderDetail{
			DeliveryOrderId: order.Id,
			CylinderId:      cyl.Id,
			BarcodeSN:       cyl.BarcodeSN,
			WeightKg:        weight,
			MasterItemId:    cyl.ItemId,
			UnitPrice:       unitPrice,
			PriceSource:     priceSource,
			LineTotal:       unitPrice,
		})
		cylinderIds = append(cylinderIds, cyl.Id)
	}

	if err := u.deliveryOrderRepo.CreateDetails(tx, details); err != nil {
		tx.Rollback()
		return nil, err
	}

	repository.LogCylinderStatusChanges(u.ledgerRepo, tx, cylinders, enum.CylinderStatusInTransit, constant.LedgerActionDOIssue, constant.AuditObjectDeliveryOrder, order.Id)
	if err := u.cylinderRepo.UpdateStatusByIds(tx, cylinderIds, enum.CylinderStatusInTransit); err != nil {
		tx.Rollback()
		return nil, err
	}

	if salesOrder != nil {
		salesOrder.Status = enum.SalesOrderStatusCompleted
		for _, line := range salesOrder.Lines {
			if line.QtyDelivered < line.QtyOrdered {
				salesOrder.Status = enum.SalesOrderStatusPartial
				break
			}
		}
		if err := u.salesOrderRepo.UpdateProgress(tx, salesOrder); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditDOIssue, constant.AuditObjectDeliveryOrder, order.Id, map[string]any{
		"do_number":       order.DONumber,
		"customer_id":     order.CustomerId,
		"fleet_id":        order.FleetId,
		"total_weight_kg": order.TotalWeightKg,
		"cylinder_qty":    order.CylinderQty,
	})

	order.Customer = *customer
	order.FleetVehicle = *fleet
	order.Details = details
	return mapper.ToDeliveryOrderResponse(order, true), nil
}

func (u *deliveryOrderUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.DeliveryOrderResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	orders, total, err := u.deliveryOrderRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.DeliveryOrderResponse]{
		Items: mapper.ToDeliveryOrderResponses(orders),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *deliveryOrderUsecase) FindById(id string) (*dto.DeliveryOrderResponse, global.ErrorResponse) {
	order, err := u.deliveryOrderRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToDeliveryOrderResponse(order, true), nil
}

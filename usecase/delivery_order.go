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
}

func NewDeliveryOrderUsecase(
	txManager helper.TxManager,
	deliveryOrderRepo repository.DeliveryOrderRepository,
	cylinderRepo repository.CylinderRepository,
	ledgerRepo repository.CylinderLedgerRepository,
	customerRepo repository.CustomerRepository,
	fleetRepo repository.FleetRepository,
	auditLogRepo repository.AuditLogRepository,
) DeliveryOrderUsecase {
	return &deliveryOrderUsecase{
		txManager:         txManager,
		deliveryOrderRepo: deliveryOrderRepo,
		cylinderRepo:      cylinderRepo,
		ledgerRepo:        ledgerRepo,
		customerRepo:      customerRepo,
		fleetRepo:         fleetRepo,
		auditLogRepo:      auditLogRepo,
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
		details = append(details, model.DeliveryOrderDetail{
			DeliveryOrderId: order.Id,
			CylinderId:      cyl.Id,
			BarcodeSN:       cyl.BarcodeSN,
			WeightKg:        weight,
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

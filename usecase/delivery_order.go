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
	"strings"
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
	customerSiteRepo  repository.CustomerSiteRepository
	fleetRepo         repository.FleetRepository
	driverRepo        repository.DriverRepository
	auditLogRepo      repository.AuditLogRepository
	pricingUsecase    CustomerItemPriceUsecase
	salesOrderRepo    repository.SalesOrderRepository
	invoiceUsecase    InvoiceUsecase
}

func NewDeliveryOrderUsecase(
	txManager helper.TxManager,
	deliveryOrderRepo repository.DeliveryOrderRepository,
	cylinderRepo repository.CylinderRepository,
	ledgerRepo repository.CylinderLedgerRepository,
	customerRepo repository.CustomerRepository,
	customerSiteRepo repository.CustomerSiteRepository,
	fleetRepo repository.FleetRepository,
	driverRepo repository.DriverRepository,
	auditLogRepo repository.AuditLogRepository,
	pricingUsecase CustomerItemPriceUsecase,
	salesOrderRepo repository.SalesOrderRepository,
	invoiceUsecase InvoiceUsecase,
) DeliveryOrderUsecase {
	return &deliveryOrderUsecase{
		txManager:         txManager,
		deliveryOrderRepo: deliveryOrderRepo,
		cylinderRepo:      cylinderRepo,
		ledgerRepo:        ledgerRepo,
		customerRepo:      customerRepo,
		customerSiteRepo:  customerSiteRepo,
		fleetRepo:         fleetRepo,
		driverRepo:        driverRepo,
		auditLogRepo:      auditLogRepo,
		pricingUsecase:    pricingUsecase,
		salesOrderRepo:    salesOrderRepo,
		invoiceUsecase:    invoiceUsecase,
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

	customerSite, siteErr := u.resolveCustomerSite(customer, req.CustomerSiteId)
	if siteErr != nil {
		return nil, siteErr
	}

	fleet, err := u.fleetRepo.FindById(req.FleetId)
	if err != nil {
		return nil, err
	}
	if !fleet.IsActive {
		return nil, global.BadRequestError("fleet vehicle is not active")
	}

	// Driver is optional for backward compatibility with clients that do not send it.
	var driver *model.Driver
	if req.DriverId != "" {
		driver, err = u.driverRepo.FindById(req.DriverId)
		if err != nil {
			return nil, err
		}
		if !driver.IsActive {
			return nil, global.BadRequestError("driver is not active")
		}
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

	var driverId *string
	if driver != nil {
		driverId = &driver.Id
	}

	order := &model.DeliveryOrder{
		BaseModel:        model.BaseModel{CreatedBy: actorUserId},
		DONumber:         helper.GenerateDONumber(),
		SalesOrderId:     req.SalesOrderId,
		CustomerId:       customer.Id,
		ShipToName:       customer.Name,
		ShipToAddress:    customer.Address,
		ShipToPostalCode: customer.PostalCode,
		ShipToPhone:      customer.Phone,
		FleetId:          fleet.Id,
		DriverId:         driverId,
		Status:           enum.DeliveryOrderStatusInTransit,
		TotalWeightKg:    totalWeight,
		CylinderQty:      len(cylinders),
		Notes:            req.Notes,
	}
	applyShipToSite(order, customerSite)

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

	invoice, err := u.invoiceUsecase.CreateForDeliveryOrder(tx, order, customer, details)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	repository.LogCylinderStatusChanges(u.ledgerRepo, tx, cylinders, enum.CylinderStatusInTransit, constant.LedgerActionDOIssue, constant.AuditObjectDeliveryOrder, order.Id)
	if err := u.cylinderRepo.UpdateStatusByIds(tx, cylinderIds, enum.CylinderStatusInTransit); err != nil {
		tx.Rollback()
		return nil, err
	}
	// Aging ("Hari di Customer") starts when the cylinders leave the warehouse on the DO, so the
	// list shows a live day count for IN_TRANSIT too, not only after the gate exchange.
	if err := u.cylinderRepo.SetOutstandingSince(tx, cylinderIds, time.Now()); err != nil {
		tx.Rollback()
		return nil, err
	}
	fleetId := fleet.Id
	if err := u.cylinderRepo.UpdateLocationByIds(tx, cylinderIds, enum.CylinderLocationInTransit, &fleetId, fleet.PlateNumber); err != nil {
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
		"driver_id":       req.DriverId,
		"total_weight_kg": order.TotalWeightKg,
		"cylinder_qty":    order.CylinderQty,
	})

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditInvoiceCreate, constant.AuditObjectInvoice, invoice.Id, map[string]any{
		"invoice_number":    invoice.InvoiceNumber,
		"delivery_order_id": order.Id,
		"total_amount":      invoice.TotalAmount,
	})

	order.Customer = *customer
	order.FleetVehicle = *fleet
	order.Driver = driver
	order.Details = details
	return mapper.ToDeliveryOrderResponse(order, true), nil
}

// resolveCustomerSite picks the explicit ship-to site (must belong to the customer)
// or falls back to the customer's primary active site. A nil site means the DO
// uses the customer head-office address.
func (u *deliveryOrderUsecase) resolveCustomerSite(customer *model.Customer, rawSiteId string) (*model.CustomerSite, global.ErrorResponse) {
	siteId := strings.TrimSpace(rawSiteId)
	sites, err := u.customerSiteRepo.FindByCustomer(customer.Id)
	if err != nil {
		return nil, err
	}

	if siteId != "" {
		for i := range sites {
			if sites[i].Id != siteId {
				continue
			}
			if !sites[i].IsActive {
				return nil, global.BadRequestError("customer site is inactive")
			}
			return &sites[i], nil
		}
		return nil, global.BadRequestError("customer site does not belong to this customer")
	}

	for i := range sites {
		if sites[i].IsPrimary && sites[i].IsActive {
			return &sites[i], nil
		}
	}
	return nil, nil
}

// applyShipToSite overlays the chosen site on the DO snapshot, keeping the customer
// value for any site field left blank.
func applyShipToSite(order *model.DeliveryOrder, site *model.CustomerSite) {
	if site == nil {
		return
	}
	order.CustomerSiteId = &site.Id
	if site.Name != "" {
		order.ShipToName = site.Name
	}
	if site.Address != "" {
		order.ShipToAddress = site.Address
	}
	if site.PostalCode != "" {
		order.ShipToPostalCode = site.PostalCode
	}
	if site.Phone != "" {
		order.ShipToPhone = site.Phone
	}
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

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

	"github.com/google/uuid"
)

type SalesOrderUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.SalesOrderResponse], global.ErrorResponse)
	FindById(id string) (*dto.SalesOrderResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateSalesOrderRequest) (*dto.SalesOrderResponse, global.ErrorResponse)
	Confirm(actorUserId, id string) (*dto.SalesOrderResponse, global.ErrorResponse)
	Cancel(actorUserId, id string) (*dto.SalesOrderResponse, global.ErrorResponse)
}

type salesOrderUsecase struct {
	txManager      helper.TxManager
	soRepo         repository.SalesOrderRepository
	poRepo         repository.CustomerPORepository
	customerRepo   repository.CustomerRepository
	masterItemRepo repository.MasterItemRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewSalesOrderUsecase(
	txManager helper.TxManager,
	soRepo repository.SalesOrderRepository,
	poRepo repository.CustomerPORepository,
	customerRepo repository.CustomerRepository,
	masterItemRepo repository.MasterItemRepository,
	auditLogRepo repository.AuditLogRepository,
) SalesOrderUsecase {
	return &salesOrderUsecase{
		txManager: txManager, soRepo: soRepo, poRepo: poRepo, customerRepo: customerRepo,
		masterItemRepo: masterItemRepo, auditLogRepo: auditLogRepo,
	}
}

func (u *salesOrderUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.SalesOrderResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	orders, total, err := u.soRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.SalesOrderResponse]{
		Items: mapper.ToSalesOrderResponses(orders),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *salesOrderUsecase) FindById(id string) (*dto.SalesOrderResponse, global.ErrorResponse) {
	order, err := u.soRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToSalesOrderResponse(order, true), nil
}

func generateSONumber() string {
	return fmt.Sprintf("SO-%s", strings.ToUpper(uuid.New().String()[:8]))
}

func (u *salesOrderUsecase) Create(actorUserId string, req *dto.CreateSalesOrderRequest) (*dto.SalesOrderResponse, global.ErrorResponse) {
	var customer *model.Customer
	var customerPOId string
	if req.CustomerPOId != "" {
		po, err := u.poRepo.FindById(req.CustomerPOId)
		if err != nil {
			return nil, err
		}
		if po.Status != enum.CustomerPOStatusConfirmed {
			return nil, global.BadRequestError("sales order can only reference a confirmed customer PO")
		}
		customer, err = u.customerRepo.FindById(po.CustomerId)
		if err != nil {
			return nil, err
		}
		customerPOId = po.Id
	} else {
		var err global.ErrorResponse
		customer, err = u.customerRepo.FindById(req.CustomerId)
		if err != nil {
			return nil, err
		}
	}
	if !customer.IsActive {
		return nil, global.BadRequestError("customer is not active")
	}

	seenItems := make(map[string]struct{}, len(req.Lines))
	lines := make([]model.SalesOrderLine, 0, len(req.Lines))
	for _, lineReq := range req.Lines {
		if _, exists := seenItems[lineReq.MasterItemId]; exists {
			return nil, global.BadRequestError("sales order cannot contain duplicate master items")
		}
		item, err := u.masterItemRepo.FindById(lineReq.MasterItemId)
		if err != nil {
			return nil, err
		}
		seenItems[lineReq.MasterItemId] = struct{}{}
		lines = append(lines, model.SalesOrderLine{MasterItemId: item.Id, QtyOrdered: lineReq.QtyOrdered})
	}

	order := &model.SalesOrder{
		BaseModel: model.BaseModel{CreatedBy: actorUserId}, SONumber: generateSONumber(),
		CustomerPOId: customerPOId, CustomerId: customer.Id, Status: enum.SalesOrderStatusDraft, Lines: lines,
	}
	tx := u.txManager.New()
	defer tx.CheckPanic()
	if err := u.soRepo.Create(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}
	_ = u.auditLogRepo.Log(actorUserId, constant.AuditSalesOrderCreate, constant.AuditObjectSalesOrder, order.Id, map[string]any{"so_number": order.SONumber, "customer_id": order.CustomerId})
	order.Customer = *customer
	return mapper.ToSalesOrderResponse(order, true), nil
}

func (u *salesOrderUsecase) transition(actorUserId, id string, target enum.SalesOrderStatus) (*dto.SalesOrderResponse, global.ErrorResponse) {
	order, err := u.soRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	if target == enum.SalesOrderStatusConfirmed && order.Status != enum.SalesOrderStatusDraft {
		return nil, global.BadRequestError("only draft sales order can be confirmed")
	}
	if target == enum.SalesOrderStatusCancelled && order.Status != enum.SalesOrderStatusDraft && order.Status != enum.SalesOrderStatusConfirmed {
		return nil, global.BadRequestError("only draft or confirmed sales order can be cancelled")
	}
	order.Status = target
	tx := u.txManager.New()
	defer tx.CheckPanic()
	if err := u.soRepo.Update(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}
	action := constant.AuditSalesOrderConfirm
	if target == enum.SalesOrderStatusCancelled {
		action = constant.AuditSalesOrderCancel
	}
	_ = u.auditLogRepo.Log(actorUserId, action, constant.AuditObjectSalesOrder, order.Id, map[string]any{"so_number": order.SONumber, "status": target, "at": time.Now().UTC()})
	return mapper.ToSalesOrderResponse(order, true), nil
}

func (u *salesOrderUsecase) Confirm(actorUserId, id string) (*dto.SalesOrderResponse, global.ErrorResponse) {
	return u.transition(actorUserId, id, enum.SalesOrderStatusConfirmed)
}

func (u *salesOrderUsecase) Cancel(actorUserId, id string) (*dto.SalesOrderResponse, global.ErrorResponse) {
	return u.transition(actorUserId, id, enum.SalesOrderStatusCancelled)
}

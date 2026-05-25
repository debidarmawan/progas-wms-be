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

type WorkOrderUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.WorkOrderResponse], global.ErrorResponse)
	FindById(id string) (*dto.WorkOrderResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateWorkOrderRequest) (*dto.WorkOrderResponse, global.ErrorResponse)
	Complete(actorUserId, id string) (*dto.WorkOrderResponse, global.ErrorResponse)
}

type workOrderUsecase struct {
	txManager          helper.TxManager
	workOrderRepo      repository.WorkOrderRepository
	masterItemRepo     repository.MasterItemRepository
	sparepartStockRepo repository.SparepartStockRepository
	movementRepo       repository.SparepartMovementRepository
	auditLogRepo       repository.AuditLogRepository
}

func NewWorkOrderUsecase(
	txManager helper.TxManager,
	workOrderRepo repository.WorkOrderRepository,
	masterItemRepo repository.MasterItemRepository,
	sparepartStockRepo repository.SparepartStockRepository,
	movementRepo repository.SparepartMovementRepository,
	auditLogRepo repository.AuditLogRepository,
) WorkOrderUsecase {
	return &workOrderUsecase{
		txManager:          txManager,
		workOrderRepo:      workOrderRepo,
		masterItemRepo:     masterItemRepo,
		sparepartStockRepo: sparepartStockRepo,
		movementRepo:       movementRepo,
		auditLogRepo:       auditLogRepo,
	}
}

func (u *workOrderUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.WorkOrderResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	orders, total, err := u.workOrderRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.WorkOrderResponse]{
		Items: mapper.ToWorkOrderResponses(orders),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *workOrderUsecase) FindById(id string) (*dto.WorkOrderResponse, global.ErrorResponse) {
	wo, err := u.workOrderRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToWorkOrderResponse(wo, true), nil
}

func (u *workOrderUsecase) Create(actorUserId string, req *dto.CreateWorkOrderRequest) (*dto.WorkOrderResponse, global.ErrorResponse) {
	tx := u.txManager.New()
	defer tx.CheckPanic()

	lines := make([]model.WorkOrderSparepart, 0, len(req.Spareparts))
	for _, line := range req.Spareparts {
		item, err := u.masterItemRepo.FindById(line.ItemId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if item.IsSerialized {
			tx.Rollback()
			return nil, global.BadRequestError(fmt.Sprintf("item %s is serialized gas, not a spare part", item.SKU))
		}
		lines = append(lines, model.WorkOrderSparepart{
			ItemId:   line.ItemId,
			Quantity: line.Quantity,
		})
	}

	wo := &model.WorkOrder{
		BaseModel:   model.BaseModel{CreatedBy: actorUserId},
		WONumber:    helper.GenerateWorkOrderNumber(),
		Title:       req.Title,
		Description: req.Description,
		Status:      enum.WorkOrderStatusOpen,
	}

	if err := u.workOrderRepo.Create(tx, wo); err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range lines {
		lines[i].WorkOrderId = wo.Id
	}
	if err := u.workOrderRepo.CreateSpareparts(tx, lines); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditWorkOrderCreate, constant.AuditObjectWorkOrder, wo.Id, map[string]any{
		"wo_number": wo.WONumber,
		"title":     wo.Title,
	})

	wo.Spareparts = lines
	return mapper.ToWorkOrderResponse(wo, true), nil
}

func (u *workOrderUsecase) Complete(actorUserId, id string) (*dto.WorkOrderResponse, global.ErrorResponse) {
	wo, err := u.workOrderRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	if wo.Status != enum.WorkOrderStatusOpen {
		return nil, global.BadRequestError("work order is not open")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	for _, line := range wo.Spareparts {
		stock, err := u.sparepartStockRepo.FindByItemIdForUpdate(tx, line.ItemId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if stock.Quantity < line.Quantity {
			tx.Rollback()
			return nil, global.BadRequestError(fmt.Sprintf(
				"insufficient stock for item %s (have %d, need %d)",
				line.ItemId, stock.Quantity, line.Quantity,
			))
		}
		before := stock.Quantity
		after := before - line.Quantity
		if err := u.sparepartStockRepo.SetQuantity(tx, line.ItemId, after); err != nil {
			tx.Rollback()
			return nil, err
		}
		movement := &model.SparepartMovement{
			ItemId:         line.ItemId,
			MovementType:   constant.MovementWorkOrder,
			QuantityDelta:  -line.Quantity,
			QuantityBefore: before,
			QuantityAfter:  after,
			ReferenceType:  constant.AuditObjectWorkOrder,
			ReferenceId:    wo.Id,
		}
		if err := u.movementRepo.Create(tx, movement); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	wo.Status = enum.WorkOrderStatusCompleted
	if err := u.workOrderRepo.Update(tx, wo); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditWorkOrderComplete, constant.AuditObjectWorkOrder, wo.Id, map[string]any{
		"wo_number": wo.WONumber,
	})

	return mapper.ToWorkOrderResponse(wo, true), nil
}

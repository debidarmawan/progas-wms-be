package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
)

type StockOpnameUsecase interface {
	Submit(actorUserId string, req *dto.StockOpnameRequest) (*dto.StockOpnameResponse, global.ErrorResponse)
}

type stockOpnameUsecase struct {
	txManager          helper.TxManager
	masterItemRepo     repository.MasterItemRepository
	sparepartStockRepo repository.SparepartStockRepository
	movementRepo       repository.SparepartMovementRepository
	auditLogRepo       repository.AuditLogRepository
}

func NewStockOpnameUsecase(
	txManager helper.TxManager,
	masterItemRepo repository.MasterItemRepository,
	sparepartStockRepo repository.SparepartStockRepository,
	movementRepo repository.SparepartMovementRepository,
	auditLogRepo repository.AuditLogRepository,
) StockOpnameUsecase {
	return &stockOpnameUsecase{
		txManager:          txManager,
		masterItemRepo:     masterItemRepo,
		sparepartStockRepo: sparepartStockRepo,
		movementRepo:       movementRepo,
		auditLogRepo:       auditLogRepo,
	}
}

func (u *stockOpnameUsecase) Submit(actorUserId string, req *dto.StockOpnameRequest) (*dto.StockOpnameResponse, global.ErrorResponse) {
	item, err := u.masterItemRepo.FindById(req.ItemId)
	if err != nil {
		return nil, err
	}
	if item.IsSerialized {
		return nil, global.BadRequestError("stock opname only applies to spare part items")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	stock, err := u.sparepartStockRepo.FindByItemIdForUpdate(tx, req.ItemId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	before := stock.Quantity
	delta := req.ActualQuantity - before
	if err := u.sparepartStockRepo.SetQuantity(tx, req.ItemId, req.ActualQuantity); err != nil {
		tx.Rollback()
		return nil, err
	}

	movement := &model.SparepartMovement{
		ItemId:         req.ItemId,
		MovementType:   constant.MovementStockOpname,
		QuantityDelta:  delta,
		QuantityBefore: before,
		QuantityAfter:  req.ActualQuantity,
		Notes:          req.Notes,
	}
	if err := u.movementRepo.Create(tx, movement); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditStockOpname, constant.AuditObjectMasterItem, item.Id, map[string]any{
		"quantity_before": before,
		"quantity_after":  req.ActualQuantity,
		"quantity_delta":  delta,
	})

	return &dto.StockOpnameResponse{
		ItemId:         item.Id,
		ItemName:       item.Name,
		QuantityBefore: before,
		QuantityAfter:  req.ActualQuantity,
		QuantityDelta:  delta,
	}, nil
}

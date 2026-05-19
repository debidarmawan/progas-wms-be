package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type MasterItemUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.MasterItemResponse], global.ErrorResponse)
	FindById(id string) (*dto.MasterItemResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateMasterItemRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateMasterItemRequest) global.ErrorResponse
}

type masterItemUsecase struct {
	txManager        helper.TxManager
	masterItemRepo   repository.MasterItemRepository
	sparepartStockRepo repository.SparepartStockRepository
	auditLogRepo     repository.AuditLogRepository
}

func NewMasterItemUsecase(
	txManager helper.TxManager,
	masterItemRepo repository.MasterItemRepository,
	sparepartStockRepo repository.SparepartStockRepository,
	auditLogRepo repository.AuditLogRepository,
) MasterItemUsecase {
	return &masterItemUsecase{
		txManager:          txManager,
		masterItemRepo:     masterItemRepo,
		sparepartStockRepo: sparepartStockRepo,
		auditLogRepo:       auditLogRepo,
	}
}

func (u *masterItemUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.MasterItemResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	items, total, err := u.masterItemRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	responses := mapper.ToMasterItemResponses(items)
	for i := range responses {
		if !items[i].IsSerialized {
			stock, stockErr := u.sparepartStockRepo.FindByItemId(items[i].Id)
			if stockErr == nil {
				qty := stock.Quantity
				responses[i].StockQuantity = &qty
			}
		}
	}
	return &dto.PaginatedResponse[dto.MasterItemResponse]{
		Items: responses,
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *masterItemUsecase) FindById(id string) (*dto.MasterItemResponse, global.ErrorResponse) {
	item, err := u.masterItemRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	var stockQty *int
	if !item.IsSerialized {
		stock, stockErr := u.sparepartStockRepo.FindByItemId(item.Id)
		if stockErr == nil {
			stockQty = &stock.Quantity
		}
	}
	return mapper.ToMasterItemResponse(item, stockQty), nil
}

func (u *masterItemUsecase) Create(actorUserId string, req *dto.CreateMasterItemRequest) global.ErrorResponse {
	existing, err := u.masterItemRepo.FindBySKU(req.SKU)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existing != nil {
		return global.BadRequestError("SKU already exists")
	}

	if req.IsSerialized && req.GasType == "" {
		return global.BadRequestError("gas_type is required for serialized gas items")
	}
	if !req.IsSerialized && req.GasType != "" {
		return global.BadRequestError("gas_type is only applicable for serialized gas items")
	}

	item := &model.MasterItem{
		Name:          req.Name,
		SKU:           req.SKU,
		GasType:       req.GasType,
		IsSerialized:  req.IsSerialized,
		EmptyWeightKg: req.EmptyWeightKg,
		GasWeightKg:   req.GasWeightKg,
		MinStockAlert: req.MinStockAlert,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.masterItemRepo.Create(tx, item); err != nil {
		tx.Rollback()
		return err
	}

	if !req.IsSerialized {
		stock := &model.SparepartStock{ItemId: item.Id, Quantity: 0}
		if err = u.sparepartStockRepo.Create(tx, stock); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditMasterItemCreate, constant.AuditObjectMasterItem, item.Id, map[string]any{
		"sku":           item.SKU,
		"is_serialized": item.IsSerialized,
	})

	return nil
}

func (u *masterItemUsecase) Update(actorUserId, id string, req *dto.UpdateMasterItemRequest) global.ErrorResponse {
	item, err := u.masterItemRepo.FindById(id)
	if err != nil {
		return err
	}

	if item.IsSerialized && req.GasType == "" {
		return global.BadRequestError("gas_type is required for serialized gas items")
	}

	item.Name = req.Name
	item.GasType = req.GasType
	item.EmptyWeightKg = req.EmptyWeightKg
	item.GasWeightKg = req.GasWeightKg
	item.MinStockAlert = req.MinStockAlert

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.masterItemRepo.Update(tx, item); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditMasterItemUpdate, constant.AuditObjectMasterItem, item.Id, map[string]any{
		"sku": item.SKU,
	})

	return nil
}

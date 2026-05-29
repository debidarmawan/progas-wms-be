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
	CreateBulk(actorUserId string, req *dto.BulkCreateMasterItemRequest) global.ErrorResponse
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
	if err := u.validateCreateRequest(req); err != nil {
		return err
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	itemId, err := u.createMasterItemInTx(tx, req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditMasterItemCreate, constant.AuditObjectMasterItem, itemId, map[string]any{
		"sku":           req.SKU,
		"is_serialized": req.IsSerialized,
	})

	return nil
}

func (u *masterItemUsecase) CreateBulk(actorUserId string, req *dto.BulkCreateMasterItemRequest) global.ErrorResponse {
	seen := map[string]bool{}
	for _, item := range req.Items {
		if err := u.validateCreateRequest(&item); err != nil {
			return err
		}
		if seen[item.SKU] {
			return global.BadRequestError("duplicate SKU in request: " + item.SKU)
		}
		seen[item.SKU] = true
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	createdIds := make([]string, 0, len(req.Items))
	for _, itemReq := range req.Items {
		itemId, err := u.createMasterItemInTx(tx, &itemReq)
		if err != nil {
			tx.Rollback()
			return err
		}
		createdIds = append(createdIds, itemId)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	for i, itemReq := range req.Items {
		_ = u.auditLogRepo.Log(actorUserId, constant.AuditMasterItemCreate, constant.AuditObjectMasterItem, createdIds[i], map[string]any{
			"sku":           itemReq.SKU,
			"is_serialized": itemReq.IsSerialized,
		})
	}

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

func (u *masterItemUsecase) validateCreateRequest(req *dto.CreateMasterItemRequest) global.ErrorResponse {
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
	return nil
}

func (u *masterItemUsecase) createMasterItemInTx(tx helper.Tx, req *dto.CreateMasterItemRequest) (string, global.ErrorResponse) {
	item := &model.MasterItem{
		Name:          req.Name,
		SKU:           req.SKU,
		GasType:       req.GasType,
		IsSerialized:  req.IsSerialized,
		EmptyWeightKg: req.EmptyWeightKg,
		GasWeightKg:   req.GasWeightKg,
		MinStockAlert: req.MinStockAlert,
	}

	if err := u.masterItemRepo.Create(tx, item); err != nil {
		return "", err
	}

	if !req.IsSerialized {
		stock := &model.SparepartStock{ItemId: item.Id, Quantity: 0}
		if err := u.sparepartStockRepo.Create(tx, stock); err != nil {
			return "", err
		}
	}

	return item.Id, nil
}

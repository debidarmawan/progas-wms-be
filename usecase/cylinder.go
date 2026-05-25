package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"time"

	"github.com/gofiber/fiber/v3"
)

type CylinderUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CylinderResponse], global.ErrorResponse)
	FindById(id string) (*dto.CylinderResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateCylinderRequest) global.ErrorResponse
}

type cylinderUsecase struct {
	txManager      helper.TxManager
	cylinderRepo   repository.CylinderRepository
	masterItemRepo repository.MasterItemRepository
	customerRepo   repository.CustomerRepository
	vendorRepo     repository.VendorRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewCylinderUsecase(
	txManager helper.TxManager,
	cylinderRepo repository.CylinderRepository,
	masterItemRepo repository.MasterItemRepository,
	customerRepo repository.CustomerRepository,
	vendorRepo repository.VendorRepository,
	auditLogRepo repository.AuditLogRepository,
) CylinderUsecase {
	return &cylinderUsecase{
		txManager:      txManager,
		cylinderRepo:   cylinderRepo,
		masterItemRepo: masterItemRepo,
		customerRepo:   customerRepo,
		vendorRepo:     vendorRepo,
		auditLogRepo:   auditLogRepo,
	}
}

func (u *cylinderUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CylinderResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	cylinders, total, err := u.cylinderRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.CylinderResponse]{
		Items: mapper.ToCylinderResponses(cylinders),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *cylinderUsecase) FindById(id string) (*dto.CylinderResponse, global.ErrorResponse) {
	cylinder, err := u.cylinderRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToCylinderResponse(cylinder), nil
}

func (u *cylinderUsecase) Create(actorUserId string, req *dto.CreateCylinderRequest) global.ErrorResponse {
	existing, err := u.cylinderRepo.FindByBarcode(req.BarcodeSN)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existing != nil {
		return global.BadRequestError("barcode already registered")
	}

	item, err := u.masterItemRepo.FindById(req.ItemId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return global.BadRequestError("invalid item")
		}
		return err
	}
	if !item.IsSerialized {
		return global.BadRequestError("item is not serialized; cannot register cylinder")
	}

	ownership := enum.Ownership(req.OwnershipType)
	var ownerId *string
	if req.OwnerId != "" {
		ownerId = &req.OwnerId
	}
	if !helper.ValidateOwnership(ownership, ownerId) {
		return global.BadRequestError("invalid ownership: CUSTOMER and VENDOR require owner_id; COMPANY must not have owner_id")
	}
	if ownership == enum.OwnershipCustomer {
		if _, err := u.customerRepo.FindById(*ownerId); err != nil {
			if err.GetCode() == fiber.StatusNotFound {
				return global.BadRequestError("invalid customer owner_id")
			}
			return err
		}
	}
	if ownership == enum.OwnershipVendor {
		vendor, err := u.vendorRepo.FindById(*ownerId)
		if err != nil {
			if err.GetCode() == fiber.StatusNotFound {
				return global.BadRequestError("invalid vendor owner_id")
			}
			return err
		}
		if !vendor.IsActive {
			return global.BadRequestError("vendor is not active")
		}
	}

	hydrotestDate, parseErr := time.Parse(time.RFC3339, req.LastHydrotestDate)
	if parseErr != nil {
		hydrotestDate, parseErr = time.Parse("2006-01-02", req.LastHydrotestDate)
	}
	if parseErr != nil {
		return global.BadRequestError("invalid last_hydrotest_date format (use YYYY-MM-DD or RFC3339)")
	}
	if !helper.ValidateHydrotestDate(hydrotestDate) {
		return global.BadRequestError("last hydrotest date is invalid or expired")
	}

	cylinder := &model.Cylinder{
		BarcodeSN:         req.BarcodeSN,
		ItemId:            req.ItemId,
		OwnershipType:     ownership,
		OwnerId:           ownerId,
		Status:            enum.CylinderStatusEmpty,
		LastHydrotestDate: hydrotestDate,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.cylinderRepo.Create(tx, cylinder); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCylinderCreate, constant.AuditObjectCylinder, cylinder.Id, map[string]any{
		"barcode_sn":     cylinder.BarcodeSN,
		"ownership_type": cylinder.OwnershipType,
		"item_id":        cylinder.ItemId,
	})

	return nil
}

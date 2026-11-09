package usecase

import (
	"strings"

	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type UomUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.UomResponse], global.ErrorResponse)
	FindById(id string) (*dto.UomResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateUomRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateUomRequest) global.ErrorResponse
	Delete(actorUserId, id string) global.ErrorResponse
}

type uomUsecase struct {
	txManager      helper.TxManager
	uomRepo        repository.UomRepository
	masterItemRepo repository.MasterItemRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewUomUsecase(
	txManager helper.TxManager,
	uomRepo repository.UomRepository,
	masterItemRepo repository.MasterItemRepository,
	auditLogRepo repository.AuditLogRepository,
) UomUsecase {
	return &uomUsecase{
		txManager:      txManager,
		uomRepo:        uomRepo,
		masterItemRepo: masterItemRepo,
		auditLogRepo:   auditLogRepo,
	}
}

func (u *uomUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.UomResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	uoms, total, err := u.uomRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.UomResponse]{
		Items: mapper.ToUomResponses(uoms),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *uomUsecase) FindById(id string) (*dto.UomResponse, global.ErrorResponse) {
	uom, err := u.uomRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToUomResponse(uom), nil
}

func (u *uomUsecase) Create(actorUserId string, req *dto.CreateUomRequest) global.ErrorResponse {
	code := normalizeUomCode(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" {
		return global.BadRequestError("code is required")
	}
	if name == "" {
		return global.BadRequestError("name is required")
	}

	if err := u.ensureCodeAvailable(code, ""); err != nil {
		return err
	}

	uom := &model.Uom{Code: code, Name: name}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.uomRepo.Create(tx, uom); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUomCreate, constant.AuditObjectUom, uom.Id, map[string]any{
		"code": uom.Code,
	})

	return nil
}

func (u *uomUsecase) Update(actorUserId, id string, req *dto.UpdateUomRequest) global.ErrorResponse {
	uom, err := u.uomRepo.FindById(id)
	if err != nil {
		return err
	}

	code := normalizeUomCode(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" {
		return global.BadRequestError("code is required")
	}
	if name == "" {
		return global.BadRequestError("name is required")
	}

	if err := u.ensureCodeAvailable(code, uom.Id); err != nil {
		return err
	}

	uom.Code = code
	uom.Name = name

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.uomRepo.Update(tx, uom); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUomUpdate, constant.AuditObjectUom, uom.Id, map[string]any{
		"code": uom.Code,
	})

	return nil
}

func (u *uomUsecase) Delete(actorUserId, id string) global.ErrorResponse {
	if _, err := u.uomRepo.FindById(id); err != nil {
		return err
	}

	usedBy, err := u.masterItemRepo.CountByUomId(id)
	if err != nil {
		return err
	}
	if usedBy > 0 {
		return global.BadRequestError("uom is used by master items and cannot be deleted")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.uomRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUomDelete, constant.AuditObjectUom, id, nil)

	return nil
}

func (u *uomUsecase) ensureCodeAvailable(code, currentId string) global.ErrorResponse {
	existing, err := u.uomRepo.FindByCode(code)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return nil
		}
		return err
	}
	if existing.Id != currentId {
		return global.BadRequestError("code already exists")
	}
	return nil
}

func normalizeUomCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

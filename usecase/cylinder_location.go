package usecase

import (
	"strings"

	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type CylinderLocationUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CylinderLocationResponse], global.ErrorResponse)
	FindById(id string) (*dto.CylinderLocationResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateCylinderLocationRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateCylinderLocationRequest) global.ErrorResponse
	Delete(actorUserId, id string) global.ErrorResponse
}

type cylinderLocationUsecase struct {
	txManager            helper.TxManager
	cylinderLocationRepo repository.CylinderLocationRepository
	auditLogRepo         repository.AuditLogRepository
}

func NewCylinderLocationUsecase(
	txManager helper.TxManager,
	cylinderLocationRepo repository.CylinderLocationRepository,
	auditLogRepo repository.AuditLogRepository,
) CylinderLocationUsecase {
	return &cylinderLocationUsecase{
		txManager:            txManager,
		cylinderLocationRepo: cylinderLocationRepo,
		auditLogRepo:         auditLogRepo,
	}
}

func (u *cylinderLocationUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CylinderLocationResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	locations, total, err := u.cylinderLocationRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.CylinderLocationResponse]{
		Items: mapper.ToCylinderLocationResponses(locations),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *cylinderLocationUsecase) FindById(id string) (*dto.CylinderLocationResponse, global.ErrorResponse) {
	location, err := u.cylinderLocationRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToCylinderLocationResponse(location), nil
}

func (u *cylinderLocationUsecase) Create(actorUserId string, req *dto.CreateCylinderLocationRequest) global.ErrorResponse {
	location := &model.CylinderLocation{IsActive: true}
	if err := u.fill(location, req.Code, req.Name, req.Type, req.Address, req.Notes); err != nil {
		return err
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.cylinderLocationRepo.Create(tx, location); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCylinderLocationCreate, constant.AuditObjectCylinderLocation, location.Id, map[string]any{
		"code": location.Code,
		"type": location.Type,
	})

	return nil
}

func (u *cylinderLocationUsecase) Update(actorUserId, id string, req *dto.UpdateCylinderLocationRequest) global.ErrorResponse {
	location, err := u.cylinderLocationRepo.FindById(id)
	if err != nil {
		return err
	}

	if err := u.fill(location, req.Code, req.Name, req.Type, req.Address, req.Notes); err != nil {
		return err
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.cylinderLocationRepo.Update(tx, location); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCylinderLocationUpdate, constant.AuditObjectCylinderLocation, location.Id, map[string]any{
		"code":      location.Code,
		"type":      location.Type,
		"is_active": location.IsActive,
	})

	return nil
}

func (u *cylinderLocationUsecase) Delete(actorUserId, id string) global.ErrorResponse {
	if _, err := u.cylinderLocationRepo.FindById(id); err != nil {
		return err
	}

	usedBy, err := u.cylinderLocationRepo.CountCylindersByLocationId(id)
	if err != nil {
		return err
	}
	if usedBy > 0 {
		return global.BadRequestError("location is still assigned to cylinders and cannot be deleted")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.cylinderLocationRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCylinderLocationDelete, constant.AuditObjectCylinderLocation, id, nil)

	return nil
}

func (u *cylinderLocationUsecase) fill(location *model.CylinderLocation, code, name, locationType, address, notes string) global.ErrorResponse {
	code = normalizeCylinderLocationCode(code)
	name = strings.TrimSpace(name)
	if code == "" {
		return global.BadRequestError("code is required")
	}
	if name == "" {
		return global.BadRequestError("name is required")
	}

	parsedType := enum.CylinderLocationType(strings.ToUpper(strings.TrimSpace(locationType)))
	if !parsedType.IsInternalPlace() {
		return global.BadRequestError("type must be WAREHOUSE, FILLING_STATION or WORKSHOP")
	}

	if err := u.ensureCodeAvailable(code, location.Id); err != nil {
		return err
	}

	location.Code = code
	location.Name = name
	location.Type = parsedType
	location.Address = strings.TrimSpace(address)
	location.Notes = strings.TrimSpace(notes)
	return nil
}

func (u *cylinderLocationUsecase) ensureCodeAvailable(code, currentId string) global.ErrorResponse {
	existing, err := u.cylinderLocationRepo.FindByCode(code)
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

func normalizeCylinderLocationCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

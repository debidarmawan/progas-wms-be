package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
)

type DriverUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.DriverResponse], global.ErrorResponse)
	FindById(id string) (*dto.DriverResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateDriverRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateDriverRequest) global.ErrorResponse
	Delete(actorUserId, id string) global.ErrorResponse
}

type driverUsecase struct {
	txManager    helper.TxManager
	driverRepo   repository.DriverRepository
	auditLogRepo repository.AuditLogRepository
}

func NewDriverUsecase(
	txManager helper.TxManager,
	driverRepo repository.DriverRepository,
	auditLogRepo repository.AuditLogRepository,
) DriverUsecase {
	return &driverUsecase{
		txManager:    txManager,
		driverRepo:   driverRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *driverUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.DriverResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	drivers, total, err := u.driverRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.DriverResponse]{
		Items: mapper.ToDriverResponses(drivers),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *driverUsecase) FindById(id string) (*dto.DriverResponse, global.ErrorResponse) {
	driver, err := u.driverRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToDriverResponse(driver), nil
}

func (u *driverUsecase) Create(actorUserId string, req *dto.CreateDriverRequest) global.ErrorResponse {
	driver := &model.Driver{
		Name:          req.Name,
		Phone:         req.Phone,
		LicenseNumber: req.LicenseNumber,
		IsActive:      true,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.driverRepo.Create(tx, driver); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditDriverCreate, constant.AuditObjectDriver, driver.Id, map[string]any{
		"name": driver.Name,
	})

	return nil
}

func (u *driverUsecase) Update(actorUserId, id string, req *dto.UpdateDriverRequest) global.ErrorResponse {
	driver, err := u.driverRepo.FindById(id)
	if err != nil {
		return err
	}

	driver.Name = req.Name
	driver.Phone = req.Phone
	driver.LicenseNumber = req.LicenseNumber
	driver.IsActive = req.IsActive

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.driverRepo.Update(tx, driver); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditDriverUpdate, constant.AuditObjectDriver, driver.Id, map[string]any{
		"name":      driver.Name,
		"is_active": driver.IsActive,
	})

	return nil
}

func (u *driverUsecase) Delete(actorUserId, id string) global.ErrorResponse {
	if _, err := u.driverRepo.FindById(id); err != nil {
		return err
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.driverRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditDriverDelete, constant.AuditObjectDriver, id, nil)

	return nil
}

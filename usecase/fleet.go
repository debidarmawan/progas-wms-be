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

type FleetUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.FleetResponse], global.ErrorResponse)
	FindById(id string) (*dto.FleetResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateFleetRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateFleetRequest) global.ErrorResponse
}

type fleetUsecase struct {
	txManager    helper.TxManager
	fleetRepo    repository.FleetRepository
	auditLogRepo repository.AuditLogRepository
}

func NewFleetUsecase(
	txManager helper.TxManager,
	fleetRepo repository.FleetRepository,
	auditLogRepo repository.AuditLogRepository,
) FleetUsecase {
	return &fleetUsecase{
		txManager:    txManager,
		fleetRepo:    fleetRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *fleetUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.FleetResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	fleets, total, err := u.fleetRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.FleetResponse]{
		Items: mapper.ToFleetResponses(fleets),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *fleetUsecase) FindById(id string) (*dto.FleetResponse, global.ErrorResponse) {
	fleet, err := u.fleetRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToFleetResponse(fleet), nil
}

func (u *fleetUsecase) Create(actorUserId string, req *dto.CreateFleetRequest) global.ErrorResponse {
	existing, err := u.fleetRepo.FindByPlate(req.PlateNumber)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existing != nil {
		return global.BadRequestError("plate number already exists")
	}

	fleet := &model.FleetVehicle{
		PlateNumber: req.PlateNumber,
		DriverName:  req.DriverName,
		MaxWeightKg: req.MaxWeightKg,
		IsActive:    true,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.fleetRepo.Create(tx, fleet); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditFleetCreate, constant.AuditObjectFleetVehicle, fleet.Id, map[string]any{
		"plate_number": fleet.PlateNumber,
	})

	return nil
}

func (u *fleetUsecase) Update(actorUserId, id string, req *dto.UpdateFleetRequest) global.ErrorResponse {
	fleet, err := u.fleetRepo.FindById(id)
	if err != nil {
		return err
	}

	fleet.DriverName = req.DriverName
	fleet.MaxWeightKg = req.MaxWeightKg
	fleet.IsActive = req.IsActive

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.fleetRepo.Update(tx, fleet); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditFleetUpdate, constant.AuditObjectFleetVehicle, fleet.Id, map[string]any{
		"plate_number": fleet.PlateNumber,
	})

	return nil
}

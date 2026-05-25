package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type FleetRepository interface {
	FindAll(page, limit int, search string) ([]model.FleetVehicle, int64, global.ErrorResponse)
	FindById(id string) (*model.FleetVehicle, global.ErrorResponse)
	FindByPlate(plate string) (*model.FleetVehicle, global.ErrorResponse)
	Create(tx helper.Tx, fleet *model.FleetVehicle) global.ErrorResponse
	Update(tx helper.Tx, fleet *model.FleetVehicle) global.ErrorResponse
}

type fleetRepository struct {
	db *gorm.DB
}

func NewFleetRepository(db *gorm.DB) FleetRepository {
	return &fleetRepository{db: db}
}

func (r *fleetRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *fleetRepository) FindAll(page, limit int, search string) ([]model.FleetVehicle, int64, global.ErrorResponse) {
	var fleets []model.FleetVehicle
	var total int64

	query := r.db.Model(&model.FleetVehicle{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("plate_number LIKE ? OR driver_name LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("plate_number asc").Offset(offset).Limit(limit).Find(&fleets).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return fleets, total, nil
}

func (r *fleetRepository) FindById(id string) (*model.FleetVehicle, global.ErrorResponse) {
	var fleet model.FleetVehicle
	if err := r.db.Where("id = ?", id).First(&fleet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Fleet vehicle not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &fleet, nil
}

func (r *fleetRepository) FindByPlate(plate string) (*model.FleetVehicle, global.ErrorResponse) {
	var fleet model.FleetVehicle
	if err := r.db.Where("plate_number = ?", plate).First(&fleet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Fleet vehicle not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &fleet, nil
}

func (r *fleetRepository) Create(tx helper.Tx, fleet *model.FleetVehicle) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(fleet).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *fleetRepository) Update(tx helper.Tx, fleet *model.FleetVehicle) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(fleet).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

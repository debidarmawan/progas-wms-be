package repository

import (
	"errors"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type CylinderLocationRepository interface {
	FindAll(page, limit int, search string) ([]model.CylinderLocation, int64, global.ErrorResponse)
	FindById(id string) (*model.CylinderLocation, global.ErrorResponse)
	FindByCode(code string) (*model.CylinderLocation, global.ErrorResponse)
	CountCylindersByLocationId(id string) (int64, global.ErrorResponse)
	Create(tx helper.Tx, location *model.CylinderLocation) global.ErrorResponse
	Update(tx helper.Tx, location *model.CylinderLocation) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type cylinderLocationRepository struct {
	db *gorm.DB
}

func NewCylinderLocationRepository(db *gorm.DB) CylinderLocationRepository {
	return &cylinderLocationRepository{db: db}
}

func (r *cylinderLocationRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *cylinderLocationRepository) FindAll(page, limit int, search string) ([]model.CylinderLocation, int64, global.ErrorResponse) {
	var locations []model.CylinderLocation
	var total int64

	query := r.db.Model(&model.CylinderLocation{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR code LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&locations).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return locations, total, nil
}

func (r *cylinderLocationRepository) FindById(id string) (*model.CylinderLocation, global.ErrorResponse) {
	var location model.CylinderLocation
	if err := r.db.Where("id = ?", id).First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Cylinder location not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &location, nil
}

func (r *cylinderLocationRepository) FindByCode(code string) (*model.CylinderLocation, global.ErrorResponse) {
	var location model.CylinderLocation
	if err := r.db.Where("code = ?", code).First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Cylinder location not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &location, nil
}

func (r *cylinderLocationRepository) CountCylindersByLocationId(id string) (int64, global.ErrorResponse) {
	var total int64
	if err := r.db.Model(&model.Cylinder{}).
		Where("location_id = ? AND location_type IN ?", id, enum.InternalCylinderLocationTypes).
		Count(&total).Error; err != nil {
		return 0, global.InternalServerError(err)
	}
	return total, nil
}

func (r *cylinderLocationRepository) Create(tx helper.Tx, location *model.CylinderLocation) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(location).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *cylinderLocationRepository) Update(tx helper.Tx, location *model.CylinderLocation) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(location).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *cylinderLocationRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.CylinderLocation{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

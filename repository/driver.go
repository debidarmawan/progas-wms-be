package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type DriverRepository interface {
	FindAll(page, limit int, search string) ([]model.Driver, int64, global.ErrorResponse)
	FindById(id string) (*model.Driver, global.ErrorResponse)
	Create(tx helper.Tx, driver *model.Driver) global.ErrorResponse
	Update(tx helper.Tx, driver *model.Driver) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *driverRepository) FindAll(page, limit int, search string) ([]model.Driver, int64, global.ErrorResponse) {
	var drivers []model.Driver
	var total int64

	query := r.db.Model(&model.Driver{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR phone LIKE ? OR license_number LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&drivers).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return drivers, total, nil
}

func (r *driverRepository) FindById(id string) (*model.Driver, global.ErrorResponse) {
	var driver model.Driver
	if err := r.db.Where("id = ?", id).First(&driver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Driver not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &driver, nil
}

func (r *driverRepository) Create(tx helper.Tx, driver *model.Driver) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(driver).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *driverRepository) Update(tx helper.Tx, driver *model.Driver) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(driver).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *driverRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.Driver{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

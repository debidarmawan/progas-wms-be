package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type UomRepository interface {
	FindAll(page, limit int, search string) ([]model.Uom, int64, global.ErrorResponse)
	FindById(id string) (*model.Uom, global.ErrorResponse)
	FindByCode(code string) (*model.Uom, global.ErrorResponse)
	Create(tx helper.Tx, uom *model.Uom) global.ErrorResponse
	Update(tx helper.Tx, uom *model.Uom) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type uomRepository struct {
	db *gorm.DB
}

func NewUomRepository(db *gorm.DB) UomRepository {
	return &uomRepository{db: db}
}

func (r *uomRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *uomRepository) FindAll(page, limit int, search string) ([]model.Uom, int64, global.ErrorResponse) {
	var uoms []model.Uom
	var total int64

	query := r.db.Model(&model.Uom{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR code LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&uoms).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return uoms, total, nil
}

func (r *uomRepository) FindById(id string) (*model.Uom, global.ErrorResponse) {
	var uom model.Uom
	if err := r.db.Where("id = ?", id).First(&uom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("UOM not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &uom, nil
}

func (r *uomRepository) FindByCode(code string) (*model.Uom, global.ErrorResponse) {
	var uom model.Uom
	if err := r.db.Where("code = ?", code).First(&uom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("UOM not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &uom, nil
}

func (r *uomRepository) Create(tx helper.Tx, uom *model.Uom) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(uom).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *uomRepository) Update(tx helper.Tx, uom *model.Uom) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(uom).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *uomRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.Uom{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

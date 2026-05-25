package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type MasterItemRepository interface {
	FindAll(page, limit int, search string) ([]model.MasterItem, int64, global.ErrorResponse)
	FindById(id string) (*model.MasterItem, global.ErrorResponse)
	FindBySKU(sku string) (*model.MasterItem, global.ErrorResponse)
	Create(tx helper.Tx, item *model.MasterItem) global.ErrorResponse
	Update(tx helper.Tx, item *model.MasterItem) global.ErrorResponse
}

type masterItemRepository struct {
	db *gorm.DB
}

func NewMasterItemRepository(db *gorm.DB) MasterItemRepository {
	return &masterItemRepository{db: db}
}

func (r *masterItemRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *masterItemRepository) FindAll(page, limit int, search string) ([]model.MasterItem, int64, global.ErrorResponse) {
	var items []model.MasterItem
	var total int64

	query := r.db.Model(&model.MasterItem{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR sku LIKE ? OR gas_type LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return items, total, nil
}

func (r *masterItemRepository) FindById(id string) (*model.MasterItem, global.ErrorResponse) {
	var item model.MasterItem
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Master item not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &item, nil
}

func (r *masterItemRepository) FindBySKU(sku string) (*model.MasterItem, global.ErrorResponse) {
	var item model.MasterItem
	if err := r.db.Where("sku = ?", sku).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Master item not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &item, nil
}

func (r *masterItemRepository) Create(tx helper.Tx, item *model.MasterItem) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(item).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *masterItemRepository) Update(tx helper.Tx, item *model.MasterItem) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(item).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

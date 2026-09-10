package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"time"

	"gorm.io/gorm"
)

type CustomerItemPriceRepository interface {
	FindAllByCustomer(customerId string, page, limit int, search string) ([]model.CustomerItemPrice, int64, global.ErrorResponse)
	FindById(id string) (*model.CustomerItemPrice, global.ErrorResponse)
	FindActive(customerId, masterItemId string, at time.Time) (*model.CustomerItemPrice, global.ErrorResponse)
	HasOverlap(customerId, masterItemId string, from time.Time, to *time.Time, excludeId string) (bool, global.ErrorResponse)
	Create(tx helper.Tx, price *model.CustomerItemPrice) global.ErrorResponse
	Update(tx helper.Tx, price *model.CustomerItemPrice) global.ErrorResponse
	Delete(tx helper.Tx, price *model.CustomerItemPrice) global.ErrorResponse
}

type customerItemPriceRepository struct {
	db *gorm.DB
}

func NewCustomerItemPriceRepository(db *gorm.DB) CustomerItemPriceRepository {
	return &customerItemPriceRepository{db: db}
}

func (r *customerItemPriceRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *customerItemPriceRepository) FindAllByCustomer(customerId string, page, limit int, search string) ([]model.CustomerItemPrice, int64, global.ErrorResponse) {
	var prices []model.CustomerItemPrice
	var total int64
	query := r.db.Model(&model.CustomerItemPrice{}).Where("customer_id = ?", customerId)
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("master_item_id IN (SELECT id FROM master_items WHERE name LIKE ? OR sku LIKE ?)", pattern, pattern)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	offset := (page - 1) * limit
	if err := query.Order("effective_from desc").Offset(offset).Limit(limit).Find(&prices).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return prices, total, nil
}

func (r *customerItemPriceRepository) FindById(id string) (*model.CustomerItemPrice, global.ErrorResponse) {
	var price model.CustomerItemPrice
	if err := r.db.Where("id = ?", id).First(&price).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Customer item price not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &price, nil
}

func (r *customerItemPriceRepository) FindActive(customerId, masterItemId string, at time.Time) (*model.CustomerItemPrice, global.ErrorResponse) {
	var price model.CustomerItemPrice
	query := r.db.Where("customer_id = ? AND master_item_id = ?", customerId, masterItemId).
		Where("effective_from <= ?", at).
		Where("effective_to IS NULL OR effective_to > ?", at).
		Order("effective_from desc")
	if err := query.First(&price).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Active customer item price not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &price, nil
}

func (r *customerItemPriceRepository) HasOverlap(customerId, masterItemId string, from time.Time, to *time.Time, excludeId string) (bool, global.ErrorResponse) {
	query := r.db.Model(&model.CustomerItemPrice{}).
		Where("customer_id = ? AND master_item_id = ?", customerId, masterItemId).
		Where("effective_from < ?", endOrMax(to)).
		Where("effective_to IS NULL OR effective_to > ?", from)
	if excludeId != "" {
		query = query.Where("id <> ?", excludeId)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, global.InternalServerError(err)
	}
	return count > 0, nil
}

func endOrMax(value *time.Time) time.Time {
	if value != nil {
		return *value
	}
	return time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
}

func (r *customerItemPriceRepository) Create(tx helper.Tx, price *model.CustomerItemPrice) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(price).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerItemPriceRepository) Update(tx helper.Tx, price *model.CustomerItemPrice) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(price).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerItemPriceRepository) Delete(tx helper.Tx, price *model.CustomerItemPrice) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(price).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

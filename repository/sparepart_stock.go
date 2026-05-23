package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type SparepartStockRepository interface {
	FindByItemId(itemId string) (*model.SparepartStock, global.ErrorResponse)
	FindByItemIdForUpdate(tx helper.Tx, itemId string) (*model.SparepartStock, global.ErrorResponse)
	FindAllLowStock() ([]model.SparepartStock, global.ErrorResponse)
	SetQuantity(tx helper.Tx, itemId string, quantity int) global.ErrorResponse
	Create(tx helper.Tx, stock *model.SparepartStock) global.ErrorResponse
}

type sparepartStockRepository struct {
	db *gorm.DB
}

func NewSparepartStockRepository(db *gorm.DB) SparepartStockRepository {
	return &sparepartStockRepository{db: db}
}

func (r *sparepartStockRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *sparepartStockRepository) FindByItemId(itemId string) (*model.SparepartStock, global.ErrorResponse) {
	var stock model.SparepartStock
	if err := r.db.Where("item_id = ?", itemId).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Sparepart stock not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &stock, nil
}

func (r *sparepartStockRepository) FindByItemIdForUpdate(tx helper.Tx, itemId string) (*model.SparepartStock, global.ErrorResponse) {
	var stock model.SparepartStock
	if err := r.dbFromTx(tx).Where("item_id = ?", itemId).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Sparepart stock not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &stock, nil
}

func (r *sparepartStockRepository) FindAllLowStock() ([]model.SparepartStock, global.ErrorResponse) {
	var stocks []model.SparepartStock
	err := r.db.Preload("MasterItem").
		Joins("JOIN master_item ON master_item.id = sparepart_stock.item_id AND master_item.deleted_at IS NULL").
		Where("sparepart_stock.quantity <= master_item.min_stock_alert").
		Find(&stocks).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	return stocks, nil
}

func (r *sparepartStockRepository) SetQuantity(tx helper.Tx, itemId string, quantity int) global.ErrorResponse {
	result := r.dbFromTx(tx).Model(&model.SparepartStock{}).
		Where("item_id = ?", itemId).
		Update("quantity", quantity)
	if result.Error != nil {
		return global.InternalServerError(result.Error)
	}
	if result.RowsAffected == 0 {
		return global.NotFoundError("Sparepart stock not found")
	}
	return nil
}

func (r *sparepartStockRepository) Create(tx helper.Tx, stock *model.SparepartStock) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(stock).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

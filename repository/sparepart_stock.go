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

func (r *sparepartStockRepository) Create(tx helper.Tx, stock *model.SparepartStock) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(stock).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

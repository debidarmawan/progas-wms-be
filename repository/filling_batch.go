package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type FillingBatchRepository interface {
	FindAll(page, limit int, search string) ([]model.FillingBatch, int64, global.ErrorResponse)
	FindById(id string) (*model.FillingBatch, global.ErrorResponse)
	Create(tx helper.Tx, batch *model.FillingBatch) global.ErrorResponse
	CreateDetails(tx helper.Tx, details []model.FillingBatchDetail) global.ErrorResponse
}

type fillingBatchRepository struct {
	db *gorm.DB
}

func NewFillingBatchRepository(db *gorm.DB) FillingBatchRepository {
	return &fillingBatchRepository{db: db}
}

func (r *fillingBatchRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *fillingBatchRepository) FindAll(page, limit int, search string) ([]model.FillingBatch, int64, global.ErrorResponse) {
	var batches []model.FillingBatch
	var total int64

	query := r.db.Model(&model.FillingBatch{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("batch_number LIKE ? OR gas_type LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Preload("MasterItem").Order("created_at desc").Offset(offset).Limit(limit).Find(&batches).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return batches, total, nil
}

func (r *fillingBatchRepository) FindById(id string) (*model.FillingBatch, global.ErrorResponse) {
	var batch model.FillingBatch
	err := r.db.Preload("MasterItem").Preload("Details.Cylinder.MasterItem").Where("id = ?", id).First(&batch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Filling batch not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &batch, nil
}

func (r *fillingBatchRepository) Create(tx helper.Tx, batch *model.FillingBatch) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(batch).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *fillingBatchRepository) CreateDetails(tx helper.Tx, details []model.FillingBatchDetail) global.ErrorResponse {
	if len(details) == 0 {
		return nil
	}
	if err := r.dbFromTx(tx).Create(&details).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

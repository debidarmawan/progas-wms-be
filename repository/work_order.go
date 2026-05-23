package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type WorkOrderRepository interface {
	FindAll(page, limit int, search string) ([]model.WorkOrder, int64, global.ErrorResponse)
	FindById(id string) (*model.WorkOrder, global.ErrorResponse)
	Create(tx helper.Tx, wo *model.WorkOrder) global.ErrorResponse
	CreateSpareparts(tx helper.Tx, lines []model.WorkOrderSparepart) global.ErrorResponse
	Update(tx helper.Tx, wo *model.WorkOrder) global.ErrorResponse
}

type workOrderRepository struct {
	db *gorm.DB
}

func NewWorkOrderRepository(db *gorm.DB) WorkOrderRepository {
	return &workOrderRepository{db: db}
}

func (r *workOrderRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *workOrderRepository) FindAll(page, limit int, search string) ([]model.WorkOrder, int64, global.ErrorResponse) {
	var orders []model.WorkOrder
	var total int64

	query := r.db.Model(&model.WorkOrder{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("wo_number LIKE ? OR title LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return orders, total, nil
}

func (r *workOrderRepository) FindById(id string) (*model.WorkOrder, global.ErrorResponse) {
	var wo model.WorkOrder
	err := r.db.Preload("Spareparts.MasterItem").Where("id = ?", id).First(&wo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Work order not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &wo, nil
}

func (r *workOrderRepository) Create(tx helper.Tx, wo *model.WorkOrder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(wo).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *workOrderRepository) CreateSpareparts(tx helper.Tx, lines []model.WorkOrderSparepart) global.ErrorResponse {
	if len(lines) == 0 {
		return nil
	}
	if err := r.dbFromTx(tx).Create(&lines).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *workOrderRepository) Update(tx helper.Tx, wo *model.WorkOrder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(wo).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

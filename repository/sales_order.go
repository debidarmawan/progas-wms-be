package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type SalesOrderRepository interface {
	FindAll(page, limit int, search string) ([]model.SalesOrder, int64, global.ErrorResponse)
	FindById(id string) (*model.SalesOrder, global.ErrorResponse)
	Create(tx helper.Tx, order *model.SalesOrder) global.ErrorResponse
	Update(tx helper.Tx, order *model.SalesOrder) global.ErrorResponse
	UpdateProgress(tx helper.Tx, order *model.SalesOrder) global.ErrorResponse
}

type salesOrderRepository struct {
	db *gorm.DB
}

func NewSalesOrderRepository(db *gorm.DB) SalesOrderRepository {
	return &salesOrderRepository{db: db}
}

func (r *salesOrderRepository) query(withLines bool) *gorm.DB {
	query := r.db.Preload("Customer")
	if withLines {
		query = query.Preload("Lines.MasterItem")
	}
	return query
}

func (r *salesOrderRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *salesOrderRepository) FindAll(page, limit int, search string) ([]model.SalesOrder, int64, global.ErrorResponse) {
	var orders []model.SalesOrder
	var total int64
	query := r.db.Model(&model.SalesOrder{}).Joins("Customer")
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("sales_order.so_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?", pattern, pattern, pattern)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	offset := (page - 1) * limit
	findQuery := r.query(false).Model(&model.SalesOrder{}).Order("sales_order.created_at desc").Offset(offset).Limit(limit)
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		findQuery = findQuery.Joins("Customer").Where("sales_order.so_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?", pattern, pattern, pattern)
	}
	if err := findQuery.Find(&orders).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return orders, total, nil
}

func (r *salesOrderRepository) FindById(id string) (*model.SalesOrder, global.ErrorResponse) {
	var order model.SalesOrder
	if err := r.query(true).Where("sales_order.id = ?", id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Sales order not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &order, nil
}

func (r *salesOrderRepository) Create(tx helper.Tx, order *model.SalesOrder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(order).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *salesOrderRepository) Update(tx helper.Tx, order *model.SalesOrder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(order).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *salesOrderRepository) UpdateProgress(tx helper.Tx, order *model.SalesOrder) global.ErrorResponse {
	db := r.dbFromTx(tx)
	for _, line := range order.Lines {
		if err := db.Model(&model.SalesOrderLine{}).Where("id = ?", line.Id).Updates(map[string]any{
			"qty_delivered": line.QtyDelivered,
		}).Error; err != nil {
			return global.InternalServerError(err)
		}
	}
	if err := db.Model(&model.SalesOrder{}).Where("id = ?", order.Id).Update("status", order.Status).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

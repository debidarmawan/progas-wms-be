package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type DeliveryOrderRepository interface {
	FindAll(page, limit int, search string) ([]model.DeliveryOrder, int64, global.ErrorResponse)
	FindById(id string) (*model.DeliveryOrder, global.ErrorResponse)
	Create(tx helper.Tx, order *model.DeliveryOrder) global.ErrorResponse
	CreateDetails(tx helper.Tx, details []model.DeliveryOrderDetail) global.ErrorResponse
}

type deliveryOrderRepository struct {
	db *gorm.DB
}

func NewDeliveryOrderRepository(db *gorm.DB) DeliveryOrderRepository {
	return &deliveryOrderRepository{db: db}
}

func (r *deliveryOrderRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *deliveryOrderRepository) FindAll(page, limit int, search string) ([]model.DeliveryOrder, int64, global.ErrorResponse) {
	var orders []model.DeliveryOrder
	var total int64

	query := r.db.Model(&model.DeliveryOrder{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Joins("Customer").Where(
			"delivery_order.do_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?",
			pattern, pattern, pattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	findQuery := r.db.Preload("Customer").Preload("FleetVehicle").Order("delivery_order.created_at desc").Offset(offset).Limit(limit)
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		findQuery = findQuery.Joins("Customer").Where(
			"delivery_order.do_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?",
			pattern, pattern, pattern,
		)
	}

	if err := findQuery.Find(&orders).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return orders, total, nil
}

func (r *deliveryOrderRepository) FindById(id string) (*model.DeliveryOrder, global.ErrorResponse) {
	var order model.DeliveryOrder
	err := r.db.Preload("Customer").Preload("FleetVehicle").Preload("Details.Cylinder.MasterItem").
		Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Delivery order not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &order, nil
}

func (r *deliveryOrderRepository) Create(tx helper.Tx, order *model.DeliveryOrder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(order).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *deliveryOrderRepository) CreateDetails(tx helper.Tx, details []model.DeliveryOrderDetail) global.ErrorResponse {
	if len(details) == 0 {
		return nil
	}
	if err := r.dbFromTx(tx).Create(&details).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

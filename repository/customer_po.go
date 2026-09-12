package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type CustomerPORepository interface {
	FindAll(page, limit int, search string) ([]model.CustomerPO, int64, global.ErrorResponse)
	FindById(id string) (*model.CustomerPO, global.ErrorResponse)
	Create(tx helper.Tx, po *model.CustomerPO) global.ErrorResponse
	Update(tx helper.Tx, po *model.CustomerPO) global.ErrorResponse
}

type customerPORepository struct {
	db *gorm.DB
}

func NewCustomerPORepository(db *gorm.DB) CustomerPORepository {
	return &customerPORepository{db: db}
}

func (r *customerPORepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *customerPORepository) query(withLines bool) *gorm.DB {
	query := r.db.Preload("Customer")
	if withLines {
		query = query.Preload("Lines.MasterItem")
	}
	return query
}

func (r *customerPORepository) FindAll(page, limit int, search string) ([]model.CustomerPO, int64, global.ErrorResponse) {
	var pos []model.CustomerPO
	var total int64

	query := r.db.Model(&model.CustomerPO{}).Joins("Customer")
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("customer_po.po_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?", pattern, pattern, pattern)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	findQuery := r.query(false).Model(&model.CustomerPO{}).Order("customer_po.created_at desc").Offset(offset).Limit(limit)
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		findQuery = findQuery.Joins("Customer").Where("customer_po.po_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?", pattern, pattern, pattern)
	}
	if err := findQuery.Find(&pos).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return pos, total, nil
}

func (r *customerPORepository) FindById(id string) (*model.CustomerPO, global.ErrorResponse) {
	var po model.CustomerPO
	if err := r.query(true).Where("customer_po.id = ?", id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Customer PO not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &po, nil
}

func (r *customerPORepository) Create(tx helper.Tx, po *model.CustomerPO) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(po).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerPORepository) Update(tx helper.Tx, po *model.CustomerPO) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(po).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

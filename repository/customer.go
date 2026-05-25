package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindAll(page, limit int, search string) ([]model.Customer, int64, global.ErrorResponse)
	FindById(id string) (*model.Customer, global.ErrorResponse)
	FindByIdForUpdate(tx helper.Tx, id string) (*model.Customer, global.ErrorResponse)
	FindByCode(code string) (*model.Customer, global.ErrorResponse)
	Create(tx helper.Tx, customer *model.Customer) global.ErrorResponse
	Update(tx helper.Tx, customer *model.Customer) global.ErrorResponse
	AdjustOutstanding(tx helper.Tx, customerId string, delta int) global.ErrorResponse
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *customerRepository) FindAll(page, limit int, search string) ([]model.Customer, int64, global.ErrorResponse) {
	var customers []model.Customer
	var total int64

	query := r.db.Model(&model.Customer{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("code LIKE ? OR name LIKE ? OR phone LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return customers, total, nil
}

func (r *customerRepository) FindById(id string) (*model.Customer, global.ErrorResponse) {
	var customer model.Customer
	if err := r.db.Where("id = ?", id).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Customer not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &customer, nil
}

func (r *customerRepository) FindByIdForUpdate(tx helper.Tx, id string) (*model.Customer, global.ErrorResponse) {
	var customer model.Customer
	if err := r.dbFromTx(tx).Where("id = ?", id).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Customer not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &customer, nil
}

func (r *customerRepository) FindByCode(code string) (*model.Customer, global.ErrorResponse) {
	var customer model.Customer
	if err := r.db.Where("code = ?", code).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Customer not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &customer, nil
}

func (r *customerRepository) Create(tx helper.Tx, customer *model.Customer) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(customer).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerRepository) Update(tx helper.Tx, customer *model.Customer) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(customer).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerRepository) AdjustOutstanding(tx helper.Tx, customerId string, delta int) global.ErrorResponse {
	result := r.dbFromTx(tx).Model(&model.Customer{}).
		Where("id = ?", customerId).
		Update("outstanding_count", gorm.Expr("outstanding_count + ?", delta))
	if result.Error != nil {
		return global.InternalServerError(result.Error)
	}
	if result.RowsAffected == 0 {
		return global.NotFoundError("Customer not found")
	}
	return nil
}

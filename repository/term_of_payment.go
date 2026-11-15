package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type TermOfPaymentRepository interface {
	FindAll(page, limit int, search string) ([]model.TermOfPayment, int64, global.ErrorResponse)
	FindById(id string) (*model.TermOfPayment, global.ErrorResponse)
	FindByCode(code string) (*model.TermOfPayment, global.ErrorResponse)
	Create(tx helper.Tx, top *model.TermOfPayment) global.ErrorResponse
	Update(tx helper.Tx, top *model.TermOfPayment) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type termOfPaymentRepository struct {
	db *gorm.DB
}

func NewTermOfPaymentRepository(db *gorm.DB) TermOfPaymentRepository {
	return &termOfPaymentRepository{db: db}
}

func (r *termOfPaymentRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *termOfPaymentRepository) FindAll(page, limit int, search string) ([]model.TermOfPayment, int64, global.ErrorResponse) {
	var tops []model.TermOfPayment
	var total int64

	query := r.db.Model(&model.TermOfPayment{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR code LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("days asc, name asc").Offset(offset).Limit(limit).Find(&tops).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return tops, total, nil
}

func (r *termOfPaymentRepository) FindById(id string) (*model.TermOfPayment, global.ErrorResponse) {
	var top model.TermOfPayment
	if err := r.db.Where("id = ?", id).First(&top).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Term of payment not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &top, nil
}

func (r *termOfPaymentRepository) FindByCode(code string) (*model.TermOfPayment, global.ErrorResponse) {
	var top model.TermOfPayment
	if err := r.db.Where("code = ?", code).First(&top).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Term of payment not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &top, nil
}

func (r *termOfPaymentRepository) Create(tx helper.Tx, top *model.TermOfPayment) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(top).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *termOfPaymentRepository) Update(tx helper.Tx, top *model.TermOfPayment) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(top).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *termOfPaymentRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.TermOfPayment{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

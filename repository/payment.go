package repository

import (
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(tx helper.Tx, payment *model.Payment) global.ErrorResponse
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *paymentRepository) Create(tx helper.Tx, payment *model.Payment) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(payment).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

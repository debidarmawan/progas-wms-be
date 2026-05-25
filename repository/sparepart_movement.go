package repository

import (
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type SparepartMovementRepository interface {
	Create(tx helper.Tx, movement *model.SparepartMovement) global.ErrorResponse
}

type sparepartMovementRepository struct {
	db *gorm.DB
}

func NewSparepartMovementRepository(db *gorm.DB) SparepartMovementRepository {
	return &sparepartMovementRepository{db: db}
}

func (r *sparepartMovementRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *sparepartMovementRepository) Create(tx helper.Tx, movement *model.SparepartMovement) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(movement).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

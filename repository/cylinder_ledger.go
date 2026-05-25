package repository

import (
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"time"

	"gorm.io/gorm"
)

type CylinderLedgerRepository interface {
	CreateBatch(tx helper.Tx, entries []model.CylinderLedger) global.ErrorResponse
	FindByBarcode(barcode string) ([]model.CylinderLedger, global.ErrorResponse)
	FindByDateRange(from, to time.Time) ([]model.CylinderLedger, global.ErrorResponse)
}

type cylinderLedgerRepository struct {
	db *gorm.DB
}

func NewCylinderLedgerRepository(db *gorm.DB) CylinderLedgerRepository {
	return &cylinderLedgerRepository{db: db}
}

func (r *cylinderLedgerRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func LogCylinderStatusChanges(
	ledgerRepo CylinderLedgerRepository,
	tx helper.Tx,
	cylinders []model.Cylinder,
	toStatus enum.CylinderStatus,
	action, refType, refId string,
) {
	entries := make([]model.CylinderLedger, 0, len(cylinders))
	for _, cyl := range cylinders {
		entries = append(entries, model.CylinderLedger{
			CylinderId:    cyl.Id,
			BarcodeSN:     cyl.BarcodeSN,
			FromStatus:    cyl.Status,
			ToStatus:      toStatus,
			Action:        action,
			ReferenceType: refType,
			ReferenceId:   refId,
		})
	}
	_ = ledgerRepo.CreateBatch(tx, entries)
}

func (r *cylinderLedgerRepository) CreateBatch(tx helper.Tx, entries []model.CylinderLedger) global.ErrorResponse {
	if len(entries) == 0 {
		return nil
	}
	if err := r.dbFromTx(tx).Create(&entries).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *cylinderLedgerRepository) FindByBarcode(barcode string) ([]model.CylinderLedger, global.ErrorResponse) {
	var entries []model.CylinderLedger
	if err := r.db.Order("created_at asc").Where("barcode_sn = ?", barcode).Find(&entries).Error; err != nil {
		return nil, global.InternalServerError(err)
	}
	return entries, nil
}

func (r *cylinderLedgerRepository) FindByDateRange(from, to time.Time) ([]model.CylinderLedger, global.ErrorResponse) {
	var entries []model.CylinderLedger
	if err := r.db.Where("created_at >= ? AND created_at <= ?", from, to).
		Order("barcode_sn asc, created_at asc").Find(&entries).Error; err != nil {
		return nil, global.InternalServerError(err)
	}
	return entries, nil
}

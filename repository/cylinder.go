package repository

import (
	"errors"
	"fmt"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"strings"

	"gorm.io/gorm"
)

type CylinderRepository interface {
	FindAll(page, limit int, search string) ([]model.Cylinder, int64, global.ErrorResponse)
	FindById(id string) (*model.Cylinder, global.ErrorResponse)
	FindByBarcode(barcode string) (*model.Cylinder, global.ErrorResponse)
	FindByBarcodes(tx helper.Tx, barcodes []string) ([]model.Cylinder, global.ErrorResponse)
	UpdateStatusByIds(tx helper.Tx, ids []string, status enum.CylinderStatus) global.ErrorResponse
	Create(tx helper.Tx, cylinder *model.Cylinder) global.ErrorResponse
}

type cylinderRepository struct {
	db *gorm.DB
}

func NewCylinderRepository(db *gorm.DB) CylinderRepository {
	return &cylinderRepository{db: db}
}

func (r *cylinderRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *cylinderRepository) FindAll(page, limit int, search string) ([]model.Cylinder, int64, global.ErrorResponse) {
	var cylinders []model.Cylinder
	var total int64

	query := r.db.Model(&model.Cylinder{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Joins("MasterItem").Where(
			"cylinder.barcode_sn LIKE ? OR cylinder.ownership_type LIKE ? OR cylinder.status LIKE ? OR master_item.name LIKE ? OR master_item.sku LIKE ? OR master_item.gas_type LIKE ?",
			pattern, pattern, pattern, pattern, pattern, pattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	findQuery := r.db.Preload("MasterItem").Order("cylinder.created_at desc").Offset(offset).Limit(limit)
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		findQuery = findQuery.Joins("MasterItem").Where(
			"cylinder.barcode_sn LIKE ? OR cylinder.ownership_type LIKE ? OR cylinder.status LIKE ? OR master_item.name LIKE ? OR master_item.sku LIKE ? OR master_item.gas_type LIKE ?",
			pattern, pattern, pattern, pattern, pattern, pattern,
		)
	}

	if err := findQuery.Find(&cylinders).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return cylinders, total, nil
}

func (r *cylinderRepository) FindById(id string) (*model.Cylinder, global.ErrorResponse) {
	var cylinder model.Cylinder
	err := r.db.Preload("MasterItem").Where("id = ?", id).First(&cylinder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Cylinder not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &cylinder, nil
}

func (r *cylinderRepository) FindByBarcode(barcode string) (*model.Cylinder, global.ErrorResponse) {
	var cylinder model.Cylinder
	err := r.db.Where("barcode_sn = ?", barcode).First(&cylinder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Cylinder not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &cylinder, nil
}

func (r *cylinderRepository) FindByBarcodes(tx helper.Tx, barcodes []string) ([]model.Cylinder, global.ErrorResponse) {
	if len(barcodes) == 0 {
		return nil, global.BadRequestError("barcodes are required")
	}

	unique := uniqueBarcodes(barcodes)
	var cylinders []model.Cylinder
	if err := r.dbFromTx(tx).Preload("MasterItem").Where("barcode_sn IN ?", unique).Find(&cylinders).Error; err != nil {
		return nil, global.InternalServerError(err)
	}

	if len(cylinders) != len(unique) {
		found := make(map[string]bool, len(cylinders))
		for _, c := range cylinders {
			found[c.BarcodeSN] = true
		}
		var missing []string
		for _, b := range unique {
			if !found[b] {
				missing = append(missing, b)
			}
		}
		return nil, global.BadRequestError(fmt.Sprintf("barcodes not found: %s", strings.Join(missing, ", ")))
	}

	return cylinders, nil
}

func (r *cylinderRepository) UpdateStatusByIds(tx helper.Tx, ids []string, status enum.CylinderStatus) global.ErrorResponse {
	if len(ids) == 0 {
		return nil
	}
	if err := r.dbFromTx(tx).Model(&model.Cylinder{}).Where("id IN ?", ids).Update("status", status).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *cylinderRepository) Create(tx helper.Tx, cylinder *model.Cylinder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(cylinder).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func uniqueBarcodes(barcodes []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, b := range barcodes {
		b = strings.TrimSpace(b)
		if b == "" || seen[b] {
			continue
		}
		seen[b] = true
		unique = append(unique, b)
	}
	return unique
}

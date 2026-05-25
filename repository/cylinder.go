package repository

import (
	"errors"
	"fmt"
	"progas-wms-be/constant"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CylinderRepository interface {
	FindAll(page, limit int, search string) ([]model.Cylinder, int64, global.ErrorResponse)
	FindById(id string) (*model.Cylinder, global.ErrorResponse)
	FindByIdForUpdate(tx helper.Tx, id string) (*model.Cylinder, global.ErrorResponse)
	FindByBarcode(barcode string) (*model.Cylinder, global.ErrorResponse)
	FindByBarcodes(tx helper.Tx, barcodes []string) ([]model.Cylinder, global.ErrorResponse)
	FindHydrotestDue(withinDays int) ([]model.Cylinder, global.ErrorResponse)
	FindOutstandingGroupedByCustomer() (map[string][]model.Cylinder, global.ErrorResponse)
	FindByVendorId(vendorId string) ([]model.Cylinder, global.ErrorResponse)
	CountByVendorIds(vendorIds []string) (map[string]int64, global.ErrorResponse)
	CountByVendorId(vendorId string) (int64, global.ErrorResponse)
	Update(tx helper.Tx, cylinder *model.Cylinder) global.ErrorResponse
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

func (r *cylinderRepository) FindByIdForUpdate(tx helper.Tx, id string) (*model.Cylinder, global.ErrorResponse) {
	var cylinder model.Cylinder
	if err := r.dbFromTx(tx).Where("id = ?", id).First(&cylinder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Cylinder not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &cylinder, nil
}

func (r *cylinderRepository) FindHydrotestDue(withinDays int) ([]model.Cylinder, global.ErrorResponse) {
	var cylinders []model.Cylinder
	now := time.Now()
	expiredBefore := now.AddDate(-constant.HydrotestValidityYears, 0, 0)
	dueSoonBy := now.AddDate(0, 0, withinDays)
	if err := r.db.Preload("MasterItem").
		Where("last_hydrotest_date < ? OR DATE_ADD(last_hydrotest_date, INTERVAL ? YEAR) <= ?",
			expiredBefore, constant.HydrotestValidityYears, dueSoonBy).
		Order("last_hydrotest_date asc").
		Find(&cylinders).Error; err != nil {
		return nil, global.InternalServerError(err)
	}
	return cylinders, nil
}

func (r *cylinderRepository) FindByVendorId(vendorId string) ([]model.Cylinder, global.ErrorResponse) {
	var cylinders []model.Cylinder
	err := r.db.Preload("MasterItem").
		Where("ownership_type = ? AND owner_id = ?", enum.OwnershipVendor, vendorId).
		Order("barcode_sn asc").
		Find(&cylinders).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	return cylinders, nil
}

func (r *cylinderRepository) CountByVendorIds(vendorIds []string) (map[string]int64, global.ErrorResponse) {
	counts := make(map[string]int64)
	if len(vendorIds) == 0 {
		return counts, nil
	}
	type row struct {
		OwnerId string
		Count   int64
	}
	var rows []row
	err := r.db.Model(&model.Cylinder{}).
		Select("owner_id, COUNT(*) as count").
		Where("ownership_type = ? AND owner_id IN ?", enum.OwnershipVendor, vendorIds).
		Group("owner_id").
		Scan(&rows).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	for _, row := range rows {
		counts[row.OwnerId] = row.Count
	}
	return counts, nil
}

func (r *cylinderRepository) CountByVendorId(vendorId string) (int64, global.ErrorResponse) {
	var count int64
	err := r.db.Model(&model.Cylinder{}).
		Where("ownership_type = ? AND owner_id = ?", enum.OwnershipVendor, vendorId).
		Count(&count).Error
	if err != nil {
		return 0, global.InternalServerError(err)
	}
	return count, nil
}

func (r *cylinderRepository) FindOutstandingGroupedByCustomer() (map[string][]model.Cylinder, global.ErrorResponse) {
	var cylinders []model.Cylinder
	if err := r.db.Where("status = ?", enum.CylinderStatusOutstanding).
		Order("owner_id asc, barcode_sn asc").Find(&cylinders).Error; err != nil {
		return nil, global.InternalServerError(err)
	}
	grouped := make(map[string][]model.Cylinder)
	for _, cyl := range cylinders {
		key := ""
		if cyl.OwnerId != nil {
			key = *cyl.OwnerId
		}
		grouped[key] = append(grouped[key], cyl)
	}
	return grouped, nil
}

func (r *cylinderRepository) Update(tx helper.Tx, cylinder *model.Cylinder) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(cylinder).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
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

package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type VendorRepository interface {
	FindAll(page, limit int, search string) ([]model.Vendor, int64, global.ErrorResponse)
	FindById(id string) (*model.Vendor, global.ErrorResponse)
	FindByCode(code string) (*model.Vendor, global.ErrorResponse)
	Create(tx helper.Tx, vendor *model.Vendor) global.ErrorResponse
	Update(tx helper.Tx, vendor *model.Vendor) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type vendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) VendorRepository {
	return &vendorRepository{db: db}
}

func (r *vendorRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *vendorRepository) FindAll(page, limit int, search string) ([]model.Vendor, int64, global.ErrorResponse) {
	var vendors []model.Vendor
	var total int64

	query := r.db.Model(&model.Vendor{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where(
			"code LIKE ? OR name LIKE ? OR phone LIKE ? OR contact_person LIKE ? OR email LIKE ?",
			pattern, pattern, pattern, pattern, pattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&vendors).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return vendors, total, nil
}

func (r *vendorRepository) FindById(id string) (*model.Vendor, global.ErrorResponse) {
	var vendor model.Vendor
	if err := r.db.Where("id = ?", id).First(&vendor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Vendor not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &vendor, nil
}

func (r *vendorRepository) FindByCode(code string) (*model.Vendor, global.ErrorResponse) {
	var vendor model.Vendor
	if err := r.db.Where("code = ?", code).First(&vendor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Vendor not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &vendor, nil
}

func (r *vendorRepository) Create(tx helper.Tx, vendor *model.Vendor) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(vendor).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *vendorRepository) Update(tx helper.Tx, vendor *model.Vendor) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(vendor).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *vendorRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.Vendor{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

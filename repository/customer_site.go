package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type CustomerSiteRepository interface {
	FindByCustomer(customerId string) ([]model.CustomerSite, global.ErrorResponse)
	FindById(id string) (*model.CustomerSite, global.ErrorResponse)
	Create(tx helper.Tx, site *model.CustomerSite) global.ErrorResponse
	Update(tx helper.Tx, site *model.CustomerSite) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
	ClearPrimary(tx helper.Tx, customerId, exceptId string) global.ErrorResponse
}

type customerSiteRepository struct {
	db *gorm.DB
}

func NewCustomerSiteRepository(db *gorm.DB) CustomerSiteRepository {
	return &customerSiteRepository{db: db}
}

func (r *customerSiteRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *customerSiteRepository) FindByCustomer(customerId string) ([]model.CustomerSite, global.ErrorResponse) {
	var sites []model.CustomerSite
	if err := r.db.Where("customer_id = ?", customerId).Order("is_primary desc, name asc").Find(&sites).Error; err != nil {
		return nil, global.InternalServerError(err)
	}
	return sites, nil
}

func (r *customerSiteRepository) FindById(id string) (*model.CustomerSite, global.ErrorResponse) {
	var site model.CustomerSite
	if err := r.db.Where("id = ?", id).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Customer site not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &site, nil
}

func (r *customerSiteRepository) Create(tx helper.Tx, site *model.CustomerSite) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(site).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerSiteRepository) Update(tx helper.Tx, site *model.CustomerSite) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(site).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerSiteRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.CustomerSite{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *customerSiteRepository) ClearPrimary(tx helper.Tx, customerId, exceptId string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Model(&model.CustomerSite{}).
		Where("customer_id = ? AND id <> ? AND is_primary = ?", customerId, exceptId, true).
		Update("is_primary", false).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

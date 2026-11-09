package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	FindAll(page, limit int, search, accountType string, isActive *bool) ([]model.Account, int64, global.ErrorResponse)
	FindById(id string) (*model.Account, global.ErrorResponse)
	FindByCode(code string) (*model.Account, global.ErrorResponse)
	CountByParentId(parentId string) (int64, global.ErrorResponse)
	Create(tx helper.Tx, account *model.Account) global.ErrorResponse
	Update(tx helper.Tx, account *model.Account) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *accountRepository) FindAll(page, limit int, search, accountType string, isActive *bool) ([]model.Account, int64, global.ErrorResponse) {
	var accounts []model.Account
	var total int64

	query := r.db.Model(&model.Account{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR code LIKE ?", pattern, pattern)
	}
	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Preload("Parent").Order("code asc").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return accounts, total, nil
}

func (r *accountRepository) FindById(id string) (*model.Account, global.ErrorResponse) {
	var account model.Account
	if err := r.db.Preload("Parent").Where("id = ?", id).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("account not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &account, nil
}

func (r *accountRepository) FindByCode(code string) (*model.Account, global.ErrorResponse) {
	var account model.Account
	if err := r.db.Where("code = ?", code).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("account not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &account, nil
}

func (r *accountRepository) CountByParentId(parentId string) (int64, global.ErrorResponse) {
	var total int64
	if err := r.db.Model(&model.Account{}).Where("parent_id = ?", parentId).Count(&total).Error; err != nil {
		return 0, global.InternalServerError(err)
	}
	return total, nil
}

func (r *accountRepository) Create(tx helper.Tx, account *model.Account) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(account).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *accountRepository) Update(tx helper.Tx, account *model.Account) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(account).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *accountRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.Account{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

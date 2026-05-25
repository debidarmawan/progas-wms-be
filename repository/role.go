package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll(page, limit int, search string) ([]model.Role, int64, global.ErrorResponse)
	FindById(id string) (*model.Role, global.ErrorResponse)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindAll(page, limit int, search string) ([]model.Role, int64, global.ErrorResponse) {
	var roles []model.Role
	var total int64

	query := r.db.Model(&model.Role{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ?", pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Order("name asc").Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return roles, total, nil
}

func (r *roleRepository) FindById(id string) (*model.Role, global.ErrorResponse) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Role not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &role, nil
}

package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]model.Role, global.ErrorResponse)
	FindById(id string) (*model.Role, global.ErrorResponse)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindAll() ([]model.Role, global.ErrorResponse) {
	var roles []model.Role
	err := r.db.Find(&roles).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	return roles, nil
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

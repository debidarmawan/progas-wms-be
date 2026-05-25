package repository

import (
	"errors"
	"progas-wms-be/constant"
	"progas-wms-be/global"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type RbacRepository interface {
	IsSuperAdmin(roleId string) (bool, global.ErrorResponse)
	HasPermission(roleId, permissionKey string) (bool, global.ErrorResponse)
}

type rbacRepository struct {
	db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) RbacRepository {
	return &rbacRepository{db: db}
}

func (r *rbacRepository) IsSuperAdmin(roleId string) (bool, global.ErrorResponse) {
	var count int64
	err := r.db.Model(&model.Role{}).
		Where("id = ? AND name = ?", roleId, constant.RoleSuperadmin).
		Count(&count).Error
	if err != nil {
		return false, global.InternalServerError(err)
	}
	return count > 0, nil
}

func (r *rbacRepository) HasPermission(roleId, permissionKey string) (bool, global.ErrorResponse) {
	var count int64
	err := r.db.Table("role_key_mapping AS rkm").
		Joins("INNER JOIN role_key AS rk ON rk.id = rkm.role_key_id AND rk.deleted_at IS NULL").
		Where("rkm.role_id = ? AND rk.key = ? AND rkm.is_allow = ? AND rkm.deleted_at IS NULL", roleId, permissionKey, true).
		Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, global.InternalServerError(err)
	}
	return count > 0, nil
}

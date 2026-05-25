package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(page, limit int, search string) ([]model.User, int64, global.ErrorResponse)
	FindByEmail(email string) (*model.User, global.ErrorResponse)
	FindByEmailExceptId(email, excludeId string) (*model.User, global.ErrorResponse)
	FindById(id string) (*model.User, global.ErrorResponse)
	UpdateLastLogin(tx helper.Tx, id string) global.ErrorResponse
	Create(tx helper.Tx, user *model.User) global.ErrorResponse
	Update(tx helper.Tx, user *model.User) global.ErrorResponse
	Delete(tx helper.Tx, id string) global.ErrorResponse
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *userRepository) FindAll(page, limit int, search string) ([]model.User, int64, global.ErrorResponse) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.Preload("Role").Order("name asc").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return users, total, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, global.ErrorResponse) {
	var user model.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("User not found")
		} else {
			return nil, global.InternalServerError(err)
		}
	}
	return &user, nil
}

func (r *userRepository) FindByEmailExceptId(email, excludeId string) (*model.User, global.ErrorResponse) {
	var user model.User
	err := r.db.Where("email = ? AND id != ?", email, excludeId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("User not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &user, nil
}

func (r *userRepository) FindById(id string) (*model.User, global.ErrorResponse) {
	var user model.User
	err := r.db.Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("User not found")
		} else {
			return nil, global.InternalServerError(err)
		}
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLogin(tx helper.Tx, id string) global.ErrorResponse {
	db := r.db
	if tx != nil {
		db = tx.Get()
	}

	now := time.Now()
	err := db.Model(&model.User{}).Where("id = ?", id).Update("last_logged_in_at", &now).Error
	if err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *userRepository) Create(tx helper.Tx, user *model.User) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(user).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *userRepository) Update(tx helper.Tx, user *model.User) global.ErrorResponse {
	if err := r.dbFromTx(tx).Save(user).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *userRepository) Delete(tx helper.Tx, id string) global.ErrorResponse {
	if err := r.dbFromTx(tx).Delete(&model.User{}, "id = ?", id).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

package repository

import (
	"progas-wms-be/global"
	"progas-wms-be/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, global.ErrorResponse)
	FindById(id string) (*model.User, global.ErrorResponse)
	UpdateLastLogin(id string) global.ErrorResponse
	Create(user *model.User) global.ErrorResponse
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, global.ErrorResponse) {
	var user model.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	return &user, nil
}

func (r *userRepository) FindById(id string) (*model.User, global.ErrorResponse) {
	var user model.User
	err := r.db.Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLogin(id string) global.ErrorResponse {
	now := time.Now()
	err := r.db.Model(&model.User{}).Where("id = ?", id).Update("last_logged_in_at", &now).Error
	if err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *userRepository) Create(user *model.User) global.ErrorResponse {
	err := r.db.Create(user).Error
	if err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

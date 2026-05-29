package config

import (
	"log"
	"progas-wms-be/constant"
	"progas-wms-be/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedBootstrapAdmin creates the first Superadmin when env vars are set and no users exist.
func SeedBootstrapAdmin(db *gorm.DB) {
	email := GetEnv("BOOTSTRAP_ADMIN_EMAIL")
	password := GetEnv("BOOTSTRAP_ADMIN_PASSWORD")
	name := GetEnv("BOOTSTRAP_ADMIN_NAME")
	if email == "" || password == "" {
		return
	}
	if name == "" {
		name = "Superadmin"
	}

	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		log.Printf("Bootstrap admin: count users: %v", err)
		return
	}
	if count > 0 {
		return
	}

	var role model.Role
	if err := db.Where("name = ?", constant.RoleSuperadmin).First(&role).Error; err != nil {
		log.Printf("Bootstrap admin: role Superadmin not found — run role seed first")
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), constant.HashingCost)
	if err != nil {
		log.Printf("Bootstrap admin: hash password: %v", err)
		return
	}
	hashed := string(hashedBytes)

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: hashed,
		RoleId:   role.Id,
		IsActive: true,
	}
	if err := db.Create(user).Error; err != nil {
		log.Printf("Bootstrap admin: create user: %v", err)
		return
	}

	log.Printf("Bootstrap admin created: %s (remove BOOTSTRAP_ADMIN_* from .env after first login)", email)
}

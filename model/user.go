package model

import "time"

type User struct {
	BaseModel
	Name           string `gorm:"not null; type:varchar(128);"`
	Email          string `gorm:"not null; type:varchar(128);uniqueIndex"`
	Phone          string `gorm:"type:varchar(32);"`
	Password       string `gorm:"not null; type:varchar(128);"`
	RoleId         string `gorm:"not null"`
	Role           Role
	DriverId       *string    `gorm:"type:varchar(36);uniqueIndex"`
	Driver         *Driver    `gorm:"foreignKey:DriverId;references:Id"`
	IsActive       bool       `gorm:""`
	ActivationDate *time.Time `gorm:""`
	LastLoggedInAt *time.Time `gorm:""`
}

package model

import "progas-wms-be/enum"

// Account is one line of the Chart of Accounts (CoA).
type Account struct {
	BaseModel
	Code     string           `gorm:"not null;type:varchar(20);uniqueIndex"`
	Name     string           `gorm:"not null;type:varchar(100)"`
	Type     enum.AccountType `gorm:"not null;type:varchar(20);index"`
	ParentId string           `gorm:"type:varchar(36);index"`
	Parent   *Account         `gorm:"foreignKey:ParentId"`
	IsActive bool             `gorm:"not null;default:true"`
}

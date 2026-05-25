package model

import "time"

type Vendor struct {
	BaseModel
	Code              string     `gorm:"type:varchar(32);uniqueIndex"`
	Name              string     `gorm:"not null;type:varchar(128)"`
	ContactPerson     string     `gorm:"type:varchar(128)"`
	Phone             string     `gorm:"type:varchar(32)"`
	Email             string     `gorm:"type:varchar(128)"`
	Address           string     `gorm:"type:varchar(255)"`
	Notes             string     `gorm:"type:varchar(500)"`
	ContractStartDate *time.Time `gorm:"type:date"`
	ContractEndDate   *time.Time `gorm:"type:date"`
	IsActive          bool       `gorm:"not null;default:true"`
}

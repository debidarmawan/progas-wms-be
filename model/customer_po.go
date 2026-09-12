package model

import (
	"progas-wms-be/enum"
	"time"
)

type CustomerPO struct {
	BaseModel
	PONumber    string                `gorm:"not null;type:varchar(50);index"`
	CustomerId  string                `gorm:"not null;type:varchar(36);index"`
	Customer    Customer              `gorm:"foreignKey:CustomerId"`
	PODate      time.Time             `gorm:"not null;type:date"`
	ValidUntil  *time.Time            `gorm:"type:date"`
	DocumentURL string                `gorm:"type:varchar(255)"`
	Status      enum.CustomerPOStatus `gorm:"not null;type:varchar(20);default:DRAFT"`
	Lines       []CustomerPOLine      `gorm:"foreignKey:CustomerPOId"`
}

type CustomerPOLine struct {
	BaseModel
	CustomerPOId string     `gorm:"not null;type:varchar(36);index"`
	MasterItemId string     `gorm:"not null;type:varchar(36);index"`
	MasterItem   MasterItem `gorm:"foreignKey:MasterItemId"`
	Quantity     int        `gorm:"not null"`
}

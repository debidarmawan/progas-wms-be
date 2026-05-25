package model

import (
	"progas-wms-be/enum"
	"time"
)

type Cylinder struct {
	BaseModel
	BarcodeSN         string             `gorm:"not null;type:varchar(50);uniqueIndex"`
	ItemId            string             `gorm:"not null;type:varchar(36);index"`
	MasterItem        MasterItem         `gorm:"foreignKey:ItemId"`
	OwnershipType     enum.Ownership     `gorm:"not null;type:varchar(20);default:'COMPANY'"`
	OwnerId           *string            `gorm:"type:varchar(36);index"`
	Status            enum.CylinderStatus `gorm:"not null;type:varchar(30);default:'EMPTY'"`
	LastHydrotestDate time.Time          `gorm:"not null"`
}

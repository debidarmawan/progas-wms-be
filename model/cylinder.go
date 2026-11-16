package model

import (
	"progas-wms-be/enum"
	"time"
)

type Cylinder struct {
	BaseModel
	BarcodeSN        string              `gorm:"not null;type:varchar(50);uniqueIndex"`
	ItemId           string              `gorm:"not null;type:varchar(36);index"`
	MasterItem       MasterItem          `gorm:"foreignKey:ItemId"`
	OwnershipType    enum.Ownership      `gorm:"not null;type:varchar(20);default:'COMPANY'"`
	OwnerId          *string             `gorm:"type:varchar(36);index"`
	Status           enum.CylinderStatus `gorm:"not null;type:varchar(30);default:'EMPTY'"`
	OutstandingSince *time.Time          `gorm:"index"`
	// MaxDaysAtCustomer is the single global rule for how long this cylinder may stay at any
	// customer (0 = no limit). MasterItem.MaxDaysAtCustomer only seeds/defaults it.
	MaxDaysAtCustomer int        `gorm:"not null;default:0"`
	LastHydrotestDate time.Time  `gorm:"not null"`
	Remarks           string     `gorm:"type:varchar(255)"`
	LabelPrintedAt    *time.Time `gorm:"index"`
	// LocationType/LocationId track where the cylinder physically is. LocationId is polymorphic
	// (CylinderLocation, Customer, Vendor or FleetVehicle per LocationType), so it carries no DB
	// foreign key and no GORM relation (a relation would create one); the pairing is enforced by
	// helper.ValidateCylinderLocation. LocationName is a snapshot so lists never join.
	LocationType enum.CylinderLocationType `gorm:"not null;type:varchar(20);default:'WAREHOUSE';index"`
	LocationId   *string                   `gorm:"type:varchar(36);index"`
	LocationName string                    `gorm:"type:varchar(128)"`
}

package model

import "progas-wms-be/enum"

// CylinderLocation is a company-owned place where cylinders physically sit (warehouse,
// filling station, workshop). Customer/vendor/in-transit locations live on Cylinder, not here.
type CylinderLocation struct {
	BaseModel
	Code     string                    `gorm:"not null;type:varchar(32);uniqueIndex"`
	Name     string                    `gorm:"not null;type:varchar(128)"`
	Type     enum.CylinderLocationType `gorm:"not null;type:varchar(20);index"`
	Address  string                    `gorm:"type:varchar(255)"`
	Notes    string                    `gorm:"type:varchar(255)"`
	IsActive bool                      `gorm:"not null;default:true"`
}

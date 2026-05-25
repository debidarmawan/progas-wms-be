package model

import "progas-wms-be/enum"

type CylinderLedger struct {
	BaseModel
	CylinderId    string              `gorm:"not null;type:varchar(36);index"`
	BarcodeSN     string              `gorm:"not null;type:varchar(50);index"`
	FromStatus    enum.CylinderStatus `gorm:"type:varchar(30)"`
	ToStatus      enum.CylinderStatus `gorm:"not null;type:varchar(30);index"`
	Action        string              `gorm:"not null;type:varchar(64)"`
	ReferenceType string              `gorm:"type:varchar(64)"`
	ReferenceId   string              `gorm:"type:varchar(36)"`
}

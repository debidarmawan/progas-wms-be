package model

import "progas-wms-be/enum"

type FillingBatch struct {
	BaseModel
	BatchNumber string                   `gorm:"not null;type:varchar(50);uniqueIndex"`
	ItemId      string                   `gorm:"not null;type:varchar(36);index"`
	MasterItem  MasterItem               `gorm:"foreignKey:ItemId"`
	GasType     string                   `gorm:"not null;type:varchar(50)"`
	Status      enum.FillingBatchStatus  `gorm:"not null;type:varchar(20)"`
	CylinderQty int                      `gorm:"not null;default:0"`
	Notes       string                   `gorm:"type:varchar(255)"`
	Details     []FillingBatchDetail     `gorm:"foreignKey:FillingBatchId"`
}

type FillingBatchDetail struct {
	BaseModel
	FillingBatchId string   `gorm:"not null;type:varchar(36);index"`
	CylinderId     string   `gorm:"not null;type:varchar(36);index"`
	BarcodeSN      string   `gorm:"not null;type:varchar(50)"`
	Cylinder       Cylinder `gorm:"foreignKey:CylinderId"`
}

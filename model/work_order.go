package model

import "progas-wms-be/enum"

type WorkOrder struct {
	BaseModel
	WONumber    string                `gorm:"not null;type:varchar(50);uniqueIndex"`
	Title       string                `gorm:"not null;type:varchar(128)"`
	Description string                `gorm:"type:varchar(500)"`
	Status      enum.WorkOrderStatus  `gorm:"not null;type:varchar(20)"`
	Spareparts  []WorkOrderSparepart  `gorm:"foreignKey:WorkOrderId"`
}

type WorkOrderSparepart struct {
	BaseModel
	WorkOrderId string     `gorm:"not null;type:varchar(36);index"`
	ItemId      string     `gorm:"not null;type:varchar(36);index"`
	Quantity    int        `gorm:"not null"`
	MasterItem  MasterItem `gorm:"foreignKey:ItemId"`
}

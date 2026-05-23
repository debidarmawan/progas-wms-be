package model

import "progas-wms-be/enum"

type DeliveryOrder struct {
	BaseModel
	DONumber      string                    `gorm:"not null;type:varchar(50);uniqueIndex"`
	CustomerId    string                    `gorm:"not null;type:varchar(36);index"`
	Customer      Customer                  `gorm:"foreignKey:CustomerId"`
	FleetId       string                    `gorm:"not null;type:varchar(36);index"`
	FleetVehicle  FleetVehicle              `gorm:"foreignKey:FleetId"`
	Status        enum.DeliveryOrderStatus  `gorm:"not null;type:varchar(30)"`
	TotalWeightKg float64                   `gorm:"type:decimal(10,2);not null;default:0"`
	CylinderQty   int                       `gorm:"not null;default:0"`
	Notes         string                    `gorm:"type:varchar(255)"`
	Details       []DeliveryOrderDetail     `gorm:"foreignKey:DeliveryOrderId"`
}

type DeliveryOrderDetail struct {
	BaseModel
	DeliveryOrderId string   `gorm:"not null;type:varchar(36);index"`
	CylinderId      string   `gorm:"not null;type:varchar(36);index"`
	BarcodeSN       string   `gorm:"not null;type:varchar(50)"`
	WeightKg        float64  `gorm:"type:decimal(8,2);not null;default:0"`
	Cylinder        Cylinder `gorm:"foreignKey:CylinderId"`
}

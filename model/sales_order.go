package model

import "progas-wms-be/enum"

type SalesOrder struct {
	BaseModel
	SONumber     string                `gorm:"not null;type:varchar(50);uniqueIndex"`
	CustomerPOId string                `gorm:"type:varchar(36);index"`
	CustomerPO   *CustomerPO           `gorm:"foreignKey:CustomerPOId"`
	CustomerId   string                `gorm:"not null;type:varchar(36);index"`
	Customer     Customer              `gorm:"foreignKey:CustomerId"`
	Status       enum.SalesOrderStatus `gorm:"not null;type:varchar(20);default:DRAFT"`
	Lines        []SalesOrderLine      `gorm:"foreignKey:SalesOrderId"`
}

type SalesOrderLine struct {
	BaseModel
	SalesOrderId string     `gorm:"not null;type:varchar(36);index"`
	MasterItemId string     `gorm:"not null;type:varchar(36);index"`
	MasterItem   MasterItem `gorm:"foreignKey:MasterItemId"`
	QtyOrdered   int        `gorm:"not null"`
	QtyDelivered int        `gorm:"not null;default:0"`
}

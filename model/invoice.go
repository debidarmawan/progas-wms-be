package model

import (
	"progas-wms-be/enum"
	"time"
)

type Invoice struct {
	BaseModel
	InvoiceNumber   string             `gorm:"not null;type:varchar(50);uniqueIndex"`
	DeliveryOrderId string             `gorm:"not null;type:varchar(36);uniqueIndex"`
	DeliveryOrder   DeliveryOrder      `gorm:"foreignKey:DeliveryOrderId"`
	CustomerId      string             `gorm:"not null;type:varchar(36);index"`
	Customer        Customer           `gorm:"foreignKey:CustomerId"`
	InvoiceDate     time.Time          `gorm:"not null"`
	DueDate         time.Time          `gorm:"not null"`
	Status          enum.InvoiceStatus `gorm:"not null;type:varchar(20);default:UNPAID"`
	TotalAmount     float64            `gorm:"type:decimal(15,2);not null;default:0"`
	PaidAmount      float64            `gorm:"type:decimal(15,2);not null;default:0"`
	Notes           string             `gorm:"type:varchar(255)"`
	Payments        []Payment          `gorm:"foreignKey:InvoiceId"`
}

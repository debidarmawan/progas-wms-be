package model

import "time"

type Payment struct {
	BaseModel
	InvoiceId   string    `gorm:"not null;type:varchar(36);index"`
	Amount      float64   `gorm:"type:decimal(15,2);not null"`
	PaidAt      time.Time `gorm:"not null"`
	Method      string    `gorm:"type:varchar(30)"`
	ReferenceNo string    `gorm:"type:varchar(100)"`
}

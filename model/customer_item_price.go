package model

import "time"

type CustomerItemPrice struct {
	BaseModel
	CustomerId    string     `gorm:"not null;type:varchar(36);index"`
	MasterItemId  string     `gorm:"not null;type:varchar(36);index"`
	Price         float64    `gorm:"not null;type:decimal(15,2)"`
	EffectiveFrom time.Time  `gorm:"not null;index"`
	EffectiveTo   *time.Time `gorm:"index"`
}

func (CustomerItemPrice) TableName() string {
	return "customer_item_prices"
}

package model

// TermOfPayment is a payment term master, e.g. CASH (0 days) or NET 30.
type TermOfPayment struct {
	BaseModel
	Code     string `gorm:"not null;type:varchar(20);uniqueIndex"`
	Name     string `gorm:"not null;type:varchar(100)"`
	Days     int    `gorm:"not null;default:0"`
	IsActive bool   `gorm:"not null;default:true"`
}

package model

// CustomerSite is a customer branch/address used as the delivery ship-to.
// Quota and outstanding counters stay on Customer, not here.
type CustomerSite struct {
	BaseModel
	CustomerId string   `gorm:"not null;type:varchar(36);index"`
	Customer   Customer `gorm:"foreignKey:CustomerId"`
	Name       string   `gorm:"not null;type:varchar(128)"`
	Pic        string   `gorm:"type:varchar(128)"`
	Phone      string   `gorm:"type:varchar(32)"`
	Address    string   `gorm:"type:varchar(255)"`
	PostalCode string   `gorm:"type:varchar(16)"`
	Notes      string   `gorm:"type:varchar(255)"`
	IsPrimary  bool     `gorm:"not null;default:false"`
	IsActive   bool     `gorm:"not null;default:true"`
}

package model

type Customer struct {
	BaseModel
	Code              string `gorm:"type:varchar(32);uniqueIndex"`
	Name              string `gorm:"not null;type:varchar(128)"`
	Phone             string `gorm:"type:varchar(32)"`
	Address           string `gorm:"type:varchar(255)"`
	CylinderQuotaLimit int   `gorm:"not null;default:0"`
	OutstandingCount  int    `gorm:"not null;default:0"`
	IsActive          bool   `gorm:"not null;default:true"`
}

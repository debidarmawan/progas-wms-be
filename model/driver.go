package model

type Driver struct {
	BaseModel
	Name          string `gorm:"not null;type:varchar(128)"`
	Phone         string `gorm:"type:varchar(32)"`
	LicenseNumber string `gorm:"type:varchar(64)"`
	IsActive      bool   `gorm:"not null;default:true"`
}

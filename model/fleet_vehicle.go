package model

type FleetVehicle struct {
	BaseModel
	PlateNumber string  `gorm:"not null;type:varchar(20);uniqueIndex"`
	DriverName  string  `gorm:"type:varchar(128)"`
	MaxWeightKg float64 `gorm:"type:decimal(8,2);not null"`
	IsActive    bool    `gorm:"not null;default:true"`
}

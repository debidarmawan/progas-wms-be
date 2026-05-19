package model

type MasterItem struct {
	BaseModel
	Name          string  `gorm:"not null;type:varchar(100)"`
	SKU           string  `gorm:"not null;type:varchar(50);uniqueIndex"`
	GasType       string  `gorm:"type:varchar(50)"`
	IsSerialized  bool    `gorm:"not null;default:true"`
	EmptyWeightKg float64 `gorm:"type:decimal(6,2)"`
	GasWeightKg   float64 `gorm:"type:decimal(6,2)"`
	MinStockAlert int     `gorm:"default:10"`
}

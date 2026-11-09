package model

// Uom is a unit of measure used by master items, e.g. PCS, KG, M3.
type Uom struct {
	BaseModel
	Name string `gorm:"not null;type:varchar(50)"`
	Code string `gorm:"not null;type:varchar(20);uniqueIndex"`
}

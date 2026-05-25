package model

type SparepartStock struct {
	BaseModel
	ItemId   string     `gorm:"not null;type:varchar(36);uniqueIndex"`
	MasterItem MasterItem `gorm:"foreignKey:ItemId"`
	Quantity int        `gorm:"not null;default:0"`
}

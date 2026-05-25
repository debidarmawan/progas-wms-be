package model

type SparepartMovement struct {
	BaseModel
	ItemId         string `gorm:"not null;type:varchar(36);index"`
	MasterItem     MasterItem `gorm:"foreignKey:ItemId"`
	MovementType   string `gorm:"not null;type:varchar(32)"`
	QuantityDelta  int    `gorm:"not null"`
	QuantityBefore int    `gorm:"not null"`
	QuantityAfter  int    `gorm:"not null"`
	ReferenceType  string `gorm:"type:varchar(64)"`
	ReferenceId    string `gorm:"type:varchar(36)"`
	Notes          string `gorm:"type:varchar(255)"`
}

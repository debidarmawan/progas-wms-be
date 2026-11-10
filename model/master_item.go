package model

type MasterItem struct {
	BaseModel
	Name              string  `gorm:"not null;type:varchar(100)"`
	SKU               string  `gorm:"not null;type:varchar(50);uniqueIndex"`
	ItemType          string  `gorm:"not null;type:varchar(10);default:gas"`
	GasType           string  `gorm:"type:varchar(50)"`
	HnaPrice          float64 `gorm:"not null;type:decimal(15,2);default:0"`
	IsSerialized      bool    `gorm:"not null;default:true"`
	EmptyWeightKg     float64 `gorm:"type:decimal(6,2)"`
	GasWeightKg       float64 `gorm:"type:decimal(6,2)"`
	MinStockAlert     int     `gorm:"default:10"`
	MaxDaysAtCustomer int     `gorm:"type:int;default:0"`
	UomId             string  `gorm:"type:varchar(36);index"`
	Uom               *Uom    `gorm:"foreignKey:UomId"`
	// CoA mappings — optional today (no auto-posting yet); used when journal posting is wired.
	InventoryAccountId string   `gorm:"type:varchar(36);index"`
	InventoryAccount   *Account `gorm:"foreignKey:InventoryAccountId"`
	RevenueAccountId   string   `gorm:"type:varchar(36);index"`
	RevenueAccount     *Account `gorm:"foreignKey:RevenueAccountId"`
	CogsAccountId      string   `gorm:"type:varchar(36);index"`
	CogsAccount        *Account `gorm:"foreignKey:CogsAccountId"`
}

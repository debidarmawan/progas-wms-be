package model

type Role struct {
	BaseModel
	Name string `gorm:"not null; type:varchar(64);unique"`
}

type RoleKey struct {
	BaseModel
	Method    string `gorm:"not null;type:varchar(8)"`
	Path      string `gorm:"not null;type:varchar(128)"`
	Key       string `gorm:"not null;type:varchar(32)"`
	KeyAccess string `gorm:"not null;type:varchar(8)"`
}

type RoleKeyMapping struct {
	BaseModel
	RoleId    string `gorm:"not null"`
	Role      Role
	RoleKeyId string `gorm:"not null"`
	RoleKey   RoleKey
	IsAllow   bool `gorm:"not null"`
}

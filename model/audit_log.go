package model

type AuditLog struct {
	BaseModel
	UserId     string `gorm:"not null;type:varchar(36);index"`
	Action     string `gorm:"not null;type:varchar(64);index"`
	ObjectType string `gorm:"type:varchar(64);index"`
	ObjectId   string `gorm:"type:varchar(36);index"`
	Details    string `gorm:"type:text"`
}

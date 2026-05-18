package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id        string `gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.Id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			log.Fatalf("failed to generate uuid v7, %v", err)
		}
		base.Id = uuidV7.String()
	}
	return nil
}

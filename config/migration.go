package config

import (
	"log"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.RoleKey{},
		&model.RoleKeyMapping{},
	)

	if err != nil {
		log.Fatal(err)
	}
}

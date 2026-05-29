package config

import (
	"errors"
	"log"
	"progas-wms-be/constant"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	names := []string{
		constant.RoleSuperadmin,
		constant.RoleWarehouseAdmin,
		constant.RoleLogisticAdmin,
		constant.RoleManager,
	}

	for _, name := range names {
		var role model.Role
		err := db.Where("name = ?", name).First(&role).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&model.Role{Name: name}).Error; err != nil {
				log.Printf("Role seed: failed to create %q: %v", name, err)
			}
			continue
		}
		if err != nil {
			log.Printf("Role seed: lookup %q: %v", name, err)
		}
	}

	log.Println("Role seed completed")
}

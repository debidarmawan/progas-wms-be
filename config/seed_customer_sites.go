package config

import (
	"log"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

// SeedCustomerSites backfills one primary site per customer from the head-office
// address so delivery orders have a default ship-to. Idempotent: customers that
// already have any (even soft-deleted) site are skipped.
func SeedCustomerSites(db *gorm.DB) {
	var customers []model.Customer
	if err := db.Where("address <> ''").Find(&customers).Error; err != nil {
		log.Printf("Customer site seed: load customers: %v", err)
		return
	}

	created := 0
	for i := range customers {
		customer := &customers[i]

		var count int64
		if err := db.Unscoped().Model(&model.CustomerSite{}).Where("customer_id = ?", customer.Id).Count(&count).Error; err != nil {
			log.Printf("Customer site seed: count sites of %s: %v", customer.Id, err)
			continue
		}
		if count > 0 {
			continue
		}

		site := model.CustomerSite{
			CustomerId: customer.Id,
			Name:       "Kantor Pusat",
			Pic:        customer.Pic,
			Phone:      customer.Phone,
			Address:    customer.Address,
			PostalCode: customer.PostalCode,
			IsPrimary:  true,
			IsActive:   true,
		}
		if err := db.Create(&site).Error; err != nil {
			log.Printf("Customer site seed: create site for %s: %v", customer.Id, err)
			continue
		}
		created++
	}

	log.Printf("Customer site seed completed, %d site(s) created", created)
}

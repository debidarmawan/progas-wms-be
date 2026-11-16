package config

import (
	"log"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

// SeedCylinderMaxDays copies the item-level aging default onto cylinders that never received a
// per-cylinder value, so the limit keeps working after it moved from Master Item to Cylinder.
// Idempotent: only rows with max_days_at_customer = 0 and a positive item default are touched.
func SeedCylinderMaxDays(db *gorm.DB) {
	var items []model.MasterItem
	if err := db.Where("max_days_at_customer > 0").Find(&items).Error; err != nil {
		log.Printf("Cylinder max days seed: load items: %v", err)
		return
	}

	updated := int64(0)
	for i := range items {
		item := &items[i]
		result := db.Model(&model.Cylinder{}).
			Where("item_id = ? AND max_days_at_customer = 0", item.Id).
			Update("max_days_at_customer", item.MaxDaysAtCustomer)
		if result.Error != nil {
			log.Printf("Cylinder max days seed: update cylinders of %s: %v", item.SKU, result.Error)
			continue
		}
		updated += result.RowsAffected
	}

	log.Printf("Cylinder max days seed completed, %d cylinder(s) updated", updated)
}

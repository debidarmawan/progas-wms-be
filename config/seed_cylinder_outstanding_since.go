package config

import (
	"log"
	"progas-wms-be/enum"

	"gorm.io/gorm"
)

// SeedOutstandingSince backfills aging for cylinders that are already away from the warehouse but
// predate the DO-issue fix. The latest delivery order detail is the last time the cylinder left,
// which is the best available start of the current stay.
// Idempotent: only rows with outstanding_since IS NULL are touched.
func SeedOutstandingSince(db *gorm.DB) {
	result := db.Exec(`
		UPDATE `+"`Cylinder`"+` c
		JOIN (
			SELECT cylinder_id, MAX(created_at) AS since
			FROM `+"`DeliveryOrderDetail`"+`
			GROUP BY cylinder_id
		) d ON d.cylinder_id = c.id
		SET c.outstanding_since = d.since
		WHERE c.outstanding_since IS NULL AND c.status IN (?, ?)`,
		enum.CylinderStatusInTransit, enum.CylinderStatusOutstanding)
	if result.Error != nil {
		log.Printf("Outstanding since seed: %v", result.Error)
		return
	}
	log.Printf("Outstanding since seed completed, %d cylinder(s) updated", result.RowsAffected)
}

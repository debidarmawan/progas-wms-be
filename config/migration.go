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
		&model.AuditLog{},
		&model.MasterItem{},
		&model.Cylinder{},
		&model.Customer{},
		&model.SparepartStock{},
		&model.FillingBatch{},
		&model.FillingBatchDetail{},
		&model.FleetVehicle{},
		&model.DeliveryOrder{},
		&model.DeliveryOrderDetail{},
		&model.WorkOrder{},
		&model.WorkOrderSparepart{},
		&model.CylinderLedger{},
		&model.SparepartMovement{},
	)

	if err != nil {
		log.Fatal(err)
	}
}

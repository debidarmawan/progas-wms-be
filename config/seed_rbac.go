package config

import (
	"errors"
	"log"
	"progas-wms-be/constant"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type permissionSeed struct {
	Key       string
	Method    string
	Path      string
	KeyAccess string
	Roles     []string
}

func SeedRBAC(db *gorm.DB) {
	permissions := []permissionSeed{
		{
			Key:       constant.PermAuthLogout,
			Method:    "POST",
			Path:      "/api/v1/logout",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermRoleRead,
			Method:    "GET",
			Path:      "/api/v1/roles",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermUserCreate,
			Method:    "POST",
			Path:      "/api/v1/users",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
			},
		},
		{
			Key:       constant.PermMasterItemRead,
			Method:    "GET",
			Path:      "/api/v1/master-items",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermMasterItemWrite,
			Method:    "POST",
			Path:      "/api/v1/master-items",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermCylinderRead,
			Method:    "GET",
			Path:      "/api/v1/cylinders",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermCylinderWrite,
			Method:    "POST",
			Path:      "/api/v1/cylinders",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermCustomerRead,
			Method:    "GET",
			Path:      "/api/v1/customers",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermCustomerWrite,
			Method:    "POST",
			Path:      "/api/v1/customers",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
			},
		},
		{
			Key:       constant.PermInboundEmptyReceive,
			Method:    "POST",
			Path:      "/api/v1/inbound/empty-receive",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermProductionQC,
			Method:    "POST",
			Path:      "/api/v1/production/qc/pre-fill",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermFillingBatchWrite,
			Method:    "POST",
			Path:      "/api/v1/production/filling-batches",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermFillingBatchRead,
			Method:    "GET",
			Path:      "/api/v1/production/filling-batches",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
	}

	for _, perm := range permissions {
		seedPermission(db, perm)
	}

	log.Println("RBAC seed completed")
}

func seedPermission(db *gorm.DB, perm permissionSeed) {
	var roleKey model.RoleKey
	err := db.Where("`key` = ?", perm.Key).First(&roleKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		roleKey = model.RoleKey{
			Method:    perm.Method,
			Path:      perm.Path,
			Key:       perm.Key,
			KeyAccess: perm.KeyAccess,
		}
		if err := db.Create(&roleKey).Error; err != nil {
			log.Printf("RBAC seed: failed to create role_key %s: %v", perm.Key, err)
			return
		}
	} else if err != nil {
		log.Printf("RBAC seed: failed to lookup role_key %s: %v", perm.Key, err)
		return
	}

	for _, roleName := range perm.Roles {
		var role model.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			log.Printf("RBAC seed: role %q not found — skip mapping for %s", roleName, perm.Key)
			continue
		}

		var existing model.RoleKeyMapping
		err := db.Where("role_id = ? AND role_key_id = ?", role.Id, roleKey.Id).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("RBAC seed: failed to lookup mapping %s/%s: %v", roleName, perm.Key, err)
			continue
		}

		mapping := model.RoleKeyMapping{
			RoleId:    role.Id,
			RoleKeyId: roleKey.Id,
			IsAllow:   true,
		}
		if err := db.Create(&mapping).Error; err != nil {
			log.Printf("RBAC seed: failed to create mapping %s/%s: %v", roleName, perm.Key, err)
		}
	}
}

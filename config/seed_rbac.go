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
			Key:       constant.PermUserRead,
			Method:    "GET",
			Path:      "/api/v1/users",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermUserWrite,
			Method:    "POST",
			Path:      "/api/v1/users",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleManager,
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
			Key:       constant.PermVendorRead,
			Method:    "GET",
			Path:      "/api/v1/vendors",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermVendorWrite,
			Method:    "POST",
			Path:      "/api/v1/vendors",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
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
			Key:       constant.PermProductionQCPreFill,
			Method:    "POST",
			Path:      "/api/v1/production/qc/pre-fill",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermProductionQCPostFill,
			Method:    "POST",
			Path:      "/api/v1/production/qc/post-fill",
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
		{
			Key:       constant.PermFleetRead,
			Method:    "GET",
			Path:      "/api/v1/logistics/fleet",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermFleetWrite,
			Method:    "POST",
			Path:      "/api/v1/logistics/fleet",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleLogisticAdmin,
			},
		},
		{
			Key:       constant.PermDORead,
			Method:    "GET",
			Path:      "/api/v1/outbound/delivery-orders",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermDOCreate,
			Method:    "POST",
			Path:      "/api/v1/outbound/delivery-orders",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleLogisticAdmin,
			},
		},
		{
			Key:       constant.PermExchangeProcess,
			Method:    "POST",
			Path:      "/api/v1/outbound/exchange",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleLogisticAdmin,
			},
		},
		{
			Key:       constant.PermExchangeApprove,
			Method:    "POST",
			Path:      "/api/v1/outbound/exchange",
			KeyAccess: "approve",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermWorkOrderRead,
			Method:    "GET",
			Path:      "/api/v1/maintenance/work-orders",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermWorkOrderWrite,
			Method:    "POST",
			Path:      "/api/v1/maintenance/work-orders",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermInventoryStockOpname,
			Method:    "POST",
			Path:      "/api/v1/inventory/spareparts/stock-opname",
			KeyAccess: "write",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
			},
		},
		{
			Key:       constant.PermCylinderHydrotest,
			Method:    "GET",
			Path:      "/api/v1/maintenance/hydrotest/due",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermDashboardRead,
			Method:    "GET",
			Path:      "/api/v1/dashboard/summary",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleWarehouseAdmin,
				constant.RoleLogisticAdmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermReportLedger,
			Method:    "GET",
			Path:      "/api/v1/reports/stock-ledger",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermReportTurnaround,
			Method:    "GET",
			Path:      "/api/v1/reports/turnaround",
			KeyAccess: "read",
			Roles: []string{
				constant.RoleSuperadmin,
				constant.RoleManager,
			},
		},
		{
			Key:       constant.PermInventoryVirtual,
			Method:    "GET",
			Path:      "/api/v1/inventory/virtual-warehouse",
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

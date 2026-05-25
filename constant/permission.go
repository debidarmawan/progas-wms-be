package constant

// Permission keys for RBAC (role_key.key).
const (
	PermAuthLogout = "auth.logout"
	PermRoleRead   = "role.read"
	PermUserRead  = "user.read"
	PermUserWrite = "user.write"

	PermMasterItemRead  = "master_item.read"
	PermMasterItemWrite = "master_item.write"
	PermCylinderRead    = "cylinder.read"
	PermCylinderWrite   = "cylinder.write"
	PermCustomerRead    = "customer.read"
	PermCustomerWrite   = "customer.write"
	PermVendorRead      = "vendor.read"
	PermVendorWrite     = "vendor.write"

	PermInboundEmptyReceive    = "inbound.empty_receive"
	PermProductionQCPreFill    = "production.qc.pre_fill"
	PermProductionQCPostFill   = "production.qc.post_fill"
	PermFillingBatchRead      = "production.filling_batch.read"
	PermFillingBatchWrite     = "production.filling_batch.write"

	PermFleetRead        = "fleet.read"
	PermFleetWrite       = "fleet.write"
	PermDORead           = "do.read"
	PermDOCreate         = "do.create"
	PermExchangeProcess  = "exchange.process"
	PermExchangeApprove  = "exchange.approve"

	PermWorkOrderRead        = "workorder.read"
	PermWorkOrderWrite       = "workorder.write"
	PermInventoryStockOpname = "inventory.stockopname"
	PermCylinderHydrotest    = "cylinder.hydrotest"
	PermDashboardRead        = "dashboard.read"
	PermReportLedger         = "report.ledger"
	PermReportTurnaround     = "report.turnaround"
	PermInventoryVirtual     = "inventory.virtual"
)

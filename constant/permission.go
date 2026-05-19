package constant

// Permission keys for RBAC (role_key.key).
const (
	PermAuthLogout = "auth.logout"
	PermRoleRead   = "role.read"
	PermUserCreate = "user.create"

	PermMasterItemRead  = "master_item.read"
	PermMasterItemWrite = "master_item.write"
	PermCylinderRead    = "cylinder.read"
	PermCylinderWrite   = "cylinder.write"
	PermCustomerRead    = "customer.read"
	PermCustomerWrite   = "customer.write"

	PermInboundEmptyReceive   = "inbound.empty_receive"
	PermProductionQC          = "production.qc"
	PermFillingBatchRead      = "production.filling_batch.read"
	PermFillingBatchWrite     = "production.filling_batch.write"
)

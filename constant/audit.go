package constant

// Audit action codes.
const (
	AuditUserCreate       = "USER_CREATE"
	AuditUserLogin        = "USER_LOGIN"
	AuditMasterItemCreate = "MASTER_ITEM_CREATE"
	AuditMasterItemUpdate = "MASTER_ITEM_UPDATE"
	AuditCylinderCreate   = "CYLINDER_CREATE"
	AuditCustomerCreate      = "CUSTOMER_CREATE"
	AuditCustomerUpdate      = "CUSTOMER_UPDATE"
	AuditEmptyReceive        = "EMPTY_RECEIVE"
	AuditPreFillQC           = "PRE_FILL_QC"
	AuditFillingBatchSubmit  = "FILLING_BATCH_SUBMIT"
)

// Audit object types.
const (
	AuditObjectUser         = "user"
	AuditObjectMasterItem   = "master_item"
	AuditObjectCylinder     = "cylinder"
	AuditObjectCustomer     = "customer"
	AuditObjectFillingBatch = "filling_batch"
)

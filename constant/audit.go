package constant

// Audit action codes.
const (
	AuditUserCreate         = "USER_CREATE"
	AuditUserUpdate         = "USER_UPDATE"
	AuditUserDelete         = "USER_DELETE"
	AuditUserLogin          = "USER_LOGIN"
	AuditMasterItemCreate   = "MASTER_ITEM_CREATE"
	AuditMasterItemUpdate   = "MASTER_ITEM_UPDATE"
	AuditCylinderCreate     = "CYLINDER_CREATE"
	AuditCustomerCreate     = "CUSTOMER_CREATE"
	AuditCustomerUpdate     = "CUSTOMER_UPDATE"
	AuditVendorCreate       = "VENDOR_CREATE"
	AuditVendorUpdate       = "VENDOR_UPDATE"
	AuditVendorDelete       = "VENDOR_DELETE"
	AuditEmptyReceive       = "EMPTY_RECEIVE"
	AuditPreFillQC          = "PRE_FILL_QC"
	AuditPostFillQC         = "POST_FILL_QC"
	AuditFillingBatchSubmit = "FILLING_BATCH_SUBMIT"
	AuditDOIssue            = "DO_ISSUE"
	AuditExchangeComplete   = "EXCHANGE_COMPLETE"
	AuditFleetCreate        = "FLEET_CREATE"
	AuditFleetUpdate        = "FLEET_UPDATE"
	AuditWorkOrderCreate    = "WORK_ORDER_CREATE"
	AuditWorkOrderComplete  = "WORK_ORDER_COMPLETE"
	AuditStockOpname        = "STOCK_OPNAME"
	AuditHydrotestRecord    = "HYDROTEST_RECORD"
)

// Audit object types.
const (
	AuditObjectUser          = "user"
	AuditObjectMasterItem    = "master_item"
	AuditObjectCylinder      = "cylinder"
	AuditObjectCustomer      = "customer"
	AuditObjectVendor        = "vendor"
	AuditObjectFillingBatch  = "filling_batch"
	AuditObjectDeliveryOrder = "delivery_order"
	AuditObjectFleetVehicle  = "fleet_vehicle"
	AuditObjectWorkOrder     = "work_order"
)

// Cylinder ledger actions.
const (
	LedgerActionEmptyReceive     = "EMPTY_RECEIVE"
	LedgerActionPreFillQC        = "PRE_FILL_QC"
	LedgerActionFillingBatch     = "FILLING_BATCH"
	LedgerActionPostFillQC       = "POST_FILL_QC"
	LedgerActionDOIssue          = "DO_ISSUE"
	LedgerActionExchangeOut      = "EXCHANGE_OUT"
	LedgerActionExchangeIn       = "EXCHANGE_IN"
	LedgerActionHydrotest        = "HYDROTEST"
	LedgerActionCylinderCreate   = "CYLINDER_CREATE"
)

// Sparepart movement types.
const (
	MovementWorkOrder   = "WORK_ORDER"
	MovementStockOpname = "STOCK_OPNAME"
)

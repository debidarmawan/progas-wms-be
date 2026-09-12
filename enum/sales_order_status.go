package enum

type SalesOrderStatus string

const (
	SalesOrderStatusDraft     SalesOrderStatus = "DRAFT"
	SalesOrderStatusConfirmed SalesOrderStatus = "CONFIRMED"
	SalesOrderStatusPartial   SalesOrderStatus = "PARTIAL"
	SalesOrderStatusCompleted SalesOrderStatus = "COMPLETED"
	SalesOrderStatusCancelled SalesOrderStatus = "CANCELLED"
)

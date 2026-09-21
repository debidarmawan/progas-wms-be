package enum

type InvoiceStatus string

const (
	InvoiceStatusUnpaid    InvoiceStatus = "UNPAID"
	InvoiceStatusPartial   InvoiceStatus = "PARTIAL"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

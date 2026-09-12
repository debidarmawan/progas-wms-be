package enum

type CustomerPOStatus string

const (
	CustomerPOStatusDraft     CustomerPOStatus = "DRAFT"
	CustomerPOStatusConfirmed CustomerPOStatus = "CONFIRMED"
	CustomerPOStatusCancelled CustomerPOStatus = "CANCELLED"
)

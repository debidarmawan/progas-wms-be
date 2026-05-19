package enum

import "database/sql/driver"

type DeliveryOrderStatus string

const (
	DeliveryOrderStatusDraft           DeliveryOrderStatus = "DRAFT"
	DeliveryOrderStatusInTransit       DeliveryOrderStatus = "IN_TRANSIT"
	DeliveryOrderStatusCompleted       DeliveryOrderStatus = "COMPLETED"
	DeliveryOrderStatusDiscrepancyHold DeliveryOrderStatus = "DISCREPANCY_HOLD"
)

func (s *DeliveryOrderStatus) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*s = DeliveryOrderStatus(v)
	return nil
}

func (s DeliveryOrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

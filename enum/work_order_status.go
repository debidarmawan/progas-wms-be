package enum

import "database/sql/driver"

type WorkOrderStatus string

const (
	WorkOrderStatusOpen      WorkOrderStatus = "OPEN"
	WorkOrderStatusCompleted WorkOrderStatus = "COMPLETED"
	WorkOrderStatusCancelled WorkOrderStatus = "CANCELLED"
)

func (s *WorkOrderStatus) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*s = WorkOrderStatus(v)
	return nil
}

func (s WorkOrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

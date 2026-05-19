package enum

import "database/sql/driver"

type FillingBatchStatus string

const (
	FillingBatchStatusCompleted FillingBatchStatus = "COMPLETED"
	FillingBatchStatusFailed    FillingBatchStatus = "FAILED"
)

func (s *FillingBatchStatus) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*s = FillingBatchStatus(v)
	return nil
}

func (s FillingBatchStatus) Value() (driver.Value, error) {
	return string(s), nil
}

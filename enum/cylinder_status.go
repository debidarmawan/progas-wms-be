package enum

import "database/sql/driver"

type CylinderStatus string

const (
	CylinderStatusEmpty       CylinderStatus = "EMPTY"
	CylinderStatusReadyToFill CylinderStatus = "READY_TO_FILL"
	CylinderStatusFilled      CylinderStatus = "FILLED"
	CylinderStatusReady       CylinderStatus = "READY"
	CylinderStatusInTransit   CylinderStatus = "IN_TRANSIT"
	CylinderStatusOutstanding CylinderStatus = "OUTSTANDING"
	CylinderStatusMaintenance CylinderStatus = "MAINTENANCE"
	CylinderStatusLost        CylinderStatus = "LOST"
	CylinderStatusWriteOff    CylinderStatus = "WRITE_OFF"
)

func (s CylinderStatus) IsValid() bool {
	switch s {
	case CylinderStatusEmpty, CylinderStatusReadyToFill, CylinderStatusFilled, CylinderStatusReady,
		CylinderStatusInTransit, CylinderStatusOutstanding, CylinderStatusMaintenance,
		CylinderStatusLost, CylinderStatusWriteOff:
		return true
	default:
		return false
	}
}

func (s *CylinderStatus) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*s = CylinderStatus(v)
	return nil
}

func (s CylinderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

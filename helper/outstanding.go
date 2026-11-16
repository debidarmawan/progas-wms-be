package helper

import (
	"progas-wms-be/enum"
	"time"
)

func DaysAtCustomer(since *time.Time) int {
	if since == nil {
		return 0
	}
	return int(time.Since(*since).Hours() / 24)
}

func IsOverdueAtCustomer(since *time.Time, maxDays int) bool {
	if maxDays <= 0 {
		return false
	}
	return DaysAtCustomer(since) >= maxDays
}

// CanSetMaxDaysAtCustomer reports whether the aging limit may still be changed for a cylinder in
// that status. LOST/WRITE_OFF cylinders have left the lifecycle, so the limit is meaningless.
func CanSetMaxDaysAtCustomer(status enum.CylinderStatus) bool {
	return status != enum.CylinderStatusLost && status != enum.CylinderStatusWriteOff
}

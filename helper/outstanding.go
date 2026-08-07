package helper

import "time"

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

package helper

import (
	"progas-wms-be/constant"
	"time"
)

func HydrotestExpiryDate(date time.Time) time.Time {
	return date.AddDate(constant.HydrotestValidityYears, 0, 0)
}

func IsHydrotestExpired(date time.Time) bool {
	return HydrotestExpiryDate(date).Before(time.Now())
}

func IsHydrotestDueSoon(date time.Time, withinDays int) bool {
	if IsHydrotestExpired(date) {
		return true
	}
	dueBy := time.Now().AddDate(0, 0, withinDays)
	return !HydrotestExpiryDate(date).After(dueBy)
}

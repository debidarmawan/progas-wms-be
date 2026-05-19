package helper

import (
	"progas-wms-be/constant"
	"progas-wms-be/enum"
	"time"
)

func ValidateHydrotestDate(date time.Time) bool {
	now := time.Now()
	if date.After(now) {
		return false
	}
	expiry := date.AddDate(constant.HydrotestValidityYears, 0, 0)
	return !expiry.Before(now)
}

func ValidateOwnership(ownerType enum.Ownership, ownerId *string) bool {
	if !ownerType.IsValid() {
		return false
	}
	if ownerType == enum.OwnershipCustomer {
		return ownerId != nil && *ownerId != ""
	}
	return true
}

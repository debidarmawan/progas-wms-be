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
	switch ownerType {
	case enum.OwnershipCompany:
		return ownerId == nil || *ownerId == ""
	case enum.OwnershipCustomer, enum.OwnershipVendor:
		return ownerId != nil && *ownerId != ""
	default:
		return false
	}
}

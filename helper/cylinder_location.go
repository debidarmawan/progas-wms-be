package helper

import (
	"fmt"
	"progas-wms-be/enum"
	"progas-wms-be/model"
	"strings"
)

// ValidateCylinderLocation guards the pairing of a location type with its reference:
// internal places must reference an active CylinderLocation of the same type, CUSTOMER/VENDOR
// must reference the party id, and IN_TRANSIT/MAINTENANCE/UNKNOWN carry no reference at all.
func ValidateCylinderLocation(locationType enum.CylinderLocationType, refId string, ref *model.CylinderLocation) error {
	if !locationType.IsValid() {
		return fmt.Errorf("invalid cylinder location type %q", locationType)
	}
	refId = strings.TrimSpace(refId)

	if locationType.IsInternalPlace() {
		if refId == "" {
			return fmt.Errorf("location_id is required for location type %s", locationType)
		}
		if ref == nil {
			return fmt.Errorf("location_id %s does not exist", refId)
		}
		if ref.Type != locationType {
			return fmt.Errorf("location %s is %s, not %s", ref.Code, ref.Type, locationType)
		}
		if !ref.IsActive {
			return fmt.Errorf("location %s is not active", ref.Code)
		}
		return nil
	}

	switch locationType {
	case enum.CylinderLocationCustomer, enum.CylinderLocationVendor:
		if refId == "" {
			return fmt.Errorf("location_id is required for location type %s", locationType)
		}
	default:
		if refId != "" {
			return fmt.Errorf("location type %s must not have location_id", locationType)
		}
	}
	return nil
}

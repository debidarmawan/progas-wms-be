package enum

import "database/sql/driver"

type CylinderLocationType string

const (
	CylinderLocationWarehouse      CylinderLocationType = "WAREHOUSE"
	CylinderLocationFillingStation CylinderLocationType = "FILLING_STATION"
	CylinderLocationWorkshop       CylinderLocationType = "WORKSHOP"
	CylinderLocationCustomer       CylinderLocationType = "CUSTOMER"
	CylinderLocationVendor         CylinderLocationType = "VENDOR"
	CylinderLocationInTransit      CylinderLocationType = "IN_TRANSIT"
	CylinderLocationMaintenance    CylinderLocationType = "MAINTENANCE"
	CylinderLocationUnknown        CylinderLocationType = "UNKNOWN"
)

// InternalCylinderLocationTypes are the company-owned places stored in the CylinderLocation
// master. CUSTOMER/VENDOR/IN_TRANSIT/MAINTENANCE/UNKNOWN are cylinder states that reference a
// customer, vendor or fleet vehicle instead of a master row.
var InternalCylinderLocationTypes = []CylinderLocationType{
	CylinderLocationWarehouse,
	CylinderLocationFillingStation,
	CylinderLocationWorkshop,
}

func (t CylinderLocationType) IsValid() bool {
	switch t {
	case CylinderLocationWarehouse, CylinderLocationFillingStation, CylinderLocationWorkshop,
		CylinderLocationCustomer, CylinderLocationVendor, CylinderLocationInTransit,
		CylinderLocationMaintenance, CylinderLocationUnknown:
		return true
	default:
		return false
	}
}

func (t CylinderLocationType) IsInternalPlace() bool {
	for _, internal := range InternalCylinderLocationTypes {
		if t == internal {
			return true
		}
	}
	return false
}

func (t *CylinderLocationType) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*t = CylinderLocationType(v)
	return nil
}

func (t CylinderLocationType) Value() (driver.Value, error) {
	return string(t), nil
}

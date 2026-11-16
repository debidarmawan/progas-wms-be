package helper

import (
	"progas-wms-be/enum"
	"progas-wms-be/model"
	"testing"
)

func TestValidateCylinderLocation(t *testing.T) {
	activeWarehouse := &model.CylinderLocation{Code: "WH-1", Type: enum.CylinderLocationWarehouse, IsActive: true}
	inactiveWarehouse := &model.CylinderLocation{Code: "WH-2", Type: enum.CylinderLocationWarehouse, IsActive: false}
	station := &model.CylinderLocation{Code: "FS-1", Type: enum.CylinderLocationFillingStation, IsActive: true}

	cases := []struct {
		name         string
		locationType enum.CylinderLocationType
		refId        string
		ref          *model.CylinderLocation
		wantErr      bool
	}{
		{"internal place with matching active master", enum.CylinderLocationWarehouse, "wh-1", activeWarehouse, false},
		{"internal place without reference", enum.CylinderLocationWarehouse, "", nil, true},
		{"internal place with unknown reference", enum.CylinderLocationWarehouse, "missing", nil, true},
		{"internal place with inactive master", enum.CylinderLocationWarehouse, "wh-2", inactiveWarehouse, true},
		{"internal place with mismatched master type", enum.CylinderLocationWorkshop, "fs-1", station, true},
		{"customer with reference id", enum.CylinderLocationCustomer, "cust-1", nil, false},
		{"customer without reference id", enum.CylinderLocationCustomer, "", nil, true},
		{"vendor with reference id", enum.CylinderLocationVendor, "vendor-1", nil, false},
		{"vendor without reference id", enum.CylinderLocationVendor, "", nil, true},
		{"in transit without reference", enum.CylinderLocationInTransit, "", nil, false},
		{"in transit with reference", enum.CylinderLocationInTransit, "fleet-1", nil, true},
		{"maintenance without reference", enum.CylinderLocationMaintenance, "", nil, false},
		{"unknown without reference", enum.CylinderLocationUnknown, "", nil, false},
		{"invalid type", enum.CylinderLocationType("GARAGE"), "", nil, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateCylinderLocation(c.locationType, c.refId, c.ref)
			if (err != nil) != c.wantErr {
				t.Fatalf("ValidateCylinderLocation(%q, %q) error = %v, wantErr %v", c.locationType, c.refId, err, c.wantErr)
			}
		})
	}
}

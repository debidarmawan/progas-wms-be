package enum

import "testing"

func TestCylinderLocationTypeIsValid(t *testing.T) {
	cases := []struct {
		name string
		in   CylinderLocationType
		want bool
	}{
		{"warehouse", CylinderLocationWarehouse, true},
		{"filling station", CylinderLocationFillingStation, true},
		{"workshop", CylinderLocationWorkshop, true},
		{"customer", CylinderLocationCustomer, true},
		{"vendor", CylinderLocationVendor, true},
		{"in transit", CylinderLocationInTransit, true},
		{"maintenance", CylinderLocationMaintenance, true},
		{"unknown", CylinderLocationUnknown, true},
		{"empty", "", false},
		{"unknown value", "GARAGE", false},
		{"lowercase", "warehouse", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.in.IsValid(); got != c.want {
				t.Errorf("IsValid(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestCylinderLocationTypeIsInternalPlace(t *testing.T) {
	cases := []struct {
		name string
		in   CylinderLocationType
		want bool
	}{
		{"warehouse", CylinderLocationWarehouse, true},
		{"filling station", CylinderLocationFillingStation, true},
		{"workshop", CylinderLocationWorkshop, true},
		{"customer", CylinderLocationCustomer, false},
		{"vendor", CylinderLocationVendor, false},
		{"in transit", CylinderLocationInTransit, false},
		{"maintenance", CylinderLocationMaintenance, false},
		{"unknown", CylinderLocationUnknown, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.in.IsInternalPlace(); got != c.want {
				t.Errorf("IsInternalPlace(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestCylinderLocationTypeRoundTrip(t *testing.T) {
	values := []CylinderLocationType{
		CylinderLocationWarehouse,
		CylinderLocationFillingStation,
		CylinderLocationWorkshop,
		CylinderLocationCustomer,
		CylinderLocationVendor,
		CylinderLocationInTransit,
		CylinderLocationMaintenance,
		CylinderLocationUnknown,
	}

	for _, want := range values {
		t.Run(string(want), func(t *testing.T) {
			value, err := want.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if value != string(want) {
				t.Fatalf("Value() = %v, want %q", value, want)
			}

			var got CylinderLocationType
			if err := got.Scan(value); err != nil {
				t.Fatalf("Scan() error = %v", err)
			}
			if got != want {
				t.Fatalf("Scan() = %q, want %q", got, want)
			}
		})
	}
}

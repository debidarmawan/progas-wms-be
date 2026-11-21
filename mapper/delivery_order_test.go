package mapper

import (
	"testing"

	"progas-wms-be/model"
)

func TestToDeliveryOrderResponseDriver(t *testing.T) {
	cases := []struct {
		name     string
		driver   *model.Driver
		wantId   string
		wantName string
	}{
		{"no driver keeps both empty", nil, "", ""},
		{"driver is mapped", &model.Driver{BaseModel: model.BaseModel{Id: "drv-1"}, Name: "Budi"}, "drv-1", "Budi"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := ToDeliveryOrderResponse(&model.DeliveryOrder{Driver: c.driver}, false)
			if res.DriverId != c.wantId || res.DriverName != c.wantName {
				t.Errorf("driver = (%q, %q), want (%q, %q)", res.DriverId, res.DriverName, c.wantId, c.wantName)
			}
		})
	}
}

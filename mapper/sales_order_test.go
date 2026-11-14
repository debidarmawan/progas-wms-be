package mapper

import (
	"progas-wms-be/model"
	"testing"
)

func TestToSalesOrderResponsePONumber(t *testing.T) {
	cases := []struct {
		name string
		po   *model.CustomerPO
		want string
	}{
		{"no customer po", nil, ""},
		{"customer po linked", &model.CustomerPO{PONumber: "PO-2026-001"}, "PO-2026-001"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			order := &model.SalesOrder{SONumber: "SO-1", CustomerPO: c.po}
			if got := ToSalesOrderResponse(order, false).PONumber; got != c.want {
				t.Errorf("PONumber = %q, want %q", got, c.want)
			}
		})
	}
}

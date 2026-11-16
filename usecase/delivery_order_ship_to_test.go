package usecase

import (
	"testing"

	"progas-wms-be/model"
)

func TestDeliveryOrderResolveCustomerSite(t *testing.T) {
	customer := &model.Customer{BaseModel: model.BaseModel{Id: "customer-id"}, Name: "PT Maju"}
	primary := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-primary"}, CustomerId: "customer-id", Name: "Kantor Pusat", IsPrimary: true, IsActive: true}
	secondary := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-2"}, CustomerId: "customer-id", Name: "Cabang Bekasi", IsActive: true}
	inactive := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-off"}, CustomerId: "customer-id", Name: "Cabang Tutup", IsActive: false}
	foreign := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-x"}, CustomerId: "other-id", Name: "Surabaya", IsActive: true}

	cases := []struct {
		name     string
		sites    []model.CustomerSite
		siteId   string
		wantSite string
		wantCode int
	}{
		{"explicit site wins", []model.CustomerSite{primary, secondary}, "site-2", "site-2", 0},
		{"falls back to primary", []model.CustomerSite{secondary, primary}, "", "site-primary", 0},
		{"no sites keeps customer address", nil, "", "", 0},
		{"rejects inactive site", []model.CustomerSite{inactive}, "site-off", "", 400},
		{"rejects site of another customer", []model.CustomerSite{primary, secondary}, foreign.Id, "", 400},
		{"ignores surrounding spaces", []model.CustomerSite{primary, secondary}, " site-2 ", "site-2", 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			usecase := &deliveryOrderUsecase{customerSiteRepo: &fakeCustomerSiteRepo{
				byCustomer: map[string][]model.CustomerSite{"customer-id": c.sites},
			}}

			site, err := usecase.resolveCustomerSite(customer, c.siteId)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("resolveCustomerSite() error = %v, want status %d", err, c.wantCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveCustomerSite() error = %v, want nil", err.Error())
			}
			gotId := ""
			if site != nil {
				gotId = site.Id
			}
			if gotId != c.wantSite {
				t.Fatalf("resolveCustomerSite() = %q, want %q", gotId, c.wantSite)
			}
		})
	}
}

func TestApplyShipToSite(t *testing.T) {
	customer := &model.Customer{Name: "PT Maju", Address: "Jl. Pusat 1", PostalCode: "40111", Phone: "021000"}
	site := &model.CustomerSite{BaseModel: model.BaseModel{Id: "site-2"}, Name: "Cabang Bekasi", Address: "Jl. Cabang 2", Phone: "021999"}

	order := &model.DeliveryOrder{
		ShipToName:       customer.Name,
		ShipToAddress:    customer.Address,
		ShipToPostalCode: customer.PostalCode,
		ShipToPhone:      customer.Phone,
	}
	applyShipToSite(order, site)

	if order.CustomerSiteId == nil || *order.CustomerSiteId != "site-2" {
		t.Fatalf("applyShipToSite() customer_site_id = %v, want site-2", order.CustomerSiteId)
	}
	if order.ShipToName != "Cabang Bekasi" || order.ShipToAddress != "Jl. Cabang 2" || order.ShipToPhone != "021999" {
		t.Fatalf("applyShipToSite() did not overlay site values: %+v", order)
	}
	// Blank site field keeps the customer value.
	if order.ShipToPostalCode != "40111" {
		t.Fatalf("applyShipToSite() postal_code = %q, want customer fallback 40111", order.ShipToPostalCode)
	}

	unchanged := &model.DeliveryOrder{ShipToName: "PT Maju", ShipToAddress: "Jl. Pusat 1"}
	applyShipToSite(unchanged, nil)
	if unchanged.CustomerSiteId != nil || unchanged.ShipToName != "PT Maju" {
		t.Fatalf("applyShipToSite(nil) changed the order: %+v", unchanged)
	}
}

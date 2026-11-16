package usecase

import (
	"strings"
	"testing"

	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
)

type fakeCustomerSiteRepo struct {
	byCustomer map[string][]model.CustomerSite
	byId       map[string]*model.CustomerSite
	created    []model.CustomerSite
	updated    *model.CustomerSite
	deletedId  string
	clearedFor string
}

func (f *fakeCustomerSiteRepo) FindByCustomer(customerId string) ([]model.CustomerSite, global.ErrorResponse) {
	return f.byCustomer[customerId], nil
}

func (f *fakeCustomerSiteRepo) FindById(id string) (*model.CustomerSite, global.ErrorResponse) {
	if site, ok := f.byId[id]; ok {
		return site, nil
	}
	return nil, global.NotFoundError("Customer site not found")
}

func (f *fakeCustomerSiteRepo) Create(tx helper.Tx, site *model.CustomerSite) global.ErrorResponse {
	site.Id = "site-created"
	f.created = append(f.created, *site)
	return nil
}

func (f *fakeCustomerSiteRepo) Update(tx helper.Tx, site *model.CustomerSite) global.ErrorResponse {
	f.updated = site
	return nil
}

func (f *fakeCustomerSiteRepo) Delete(tx helper.Tx, id string) global.ErrorResponse {
	f.deletedId = id
	return nil
}

func (f *fakeCustomerSiteRepo) ClearPrimary(tx helper.Tx, customerId, exceptId string) global.ErrorResponse {
	f.clearedFor = customerId
	return nil
}

func newFakeCustomerSiteUsecase(customerRepo *fakeCustomerRepo, siteRepo *fakeCustomerSiteRepo) CustomerSiteUsecase {
	return NewCustomerSiteUsecase(fakeTxManager{}, customerRepo, siteRepo, fakeAuditLogRepo{})
}

func newCustomerSiteTestRepos(sites ...model.CustomerSite) (*fakeCustomerRepo, *fakeCustomerSiteRepo) {
	customer := &model.Customer{BaseModel: model.BaseModel{Id: "customer-id"}, Name: "PT Maju"}
	customerRepo := &fakeCustomerRepo{byId: map[string]*model.Customer{"customer-id": customer}}

	byId := make(map[string]*model.CustomerSite, len(sites))
	for i := range sites {
		byId[sites[i].Id] = &sites[i]
	}
	siteRepo := &fakeCustomerSiteRepo{
		byCustomer: map[string][]model.CustomerSite{"customer-id": sites},
		byId:       byId,
	}
	return customerRepo, siteRepo
}

func TestCustomerSiteUsecaseCreate(t *testing.T) {
	existingSite := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-1"}, CustomerId: "customer-id", Name: "Bandung", IsPrimary: true, IsActive: true}

	cases := []struct {
		name        string
		customerId  string
		req         dto.CreateCustomerSiteRequest
		existing    []model.CustomerSite
		wantCode    int
		wantPrimary bool
		wantCleared bool
		wantCreated int
	}{
		{
			name:        "first site of a customer becomes primary",
			customerId:  "customer-id",
			req:         dto.CreateCustomerSiteRequest{Name: " Kantor Pusat "},
			wantPrimary: true,
			wantCreated: 1,
		},
		{
			name:        "second site is not primary unless requested",
			customerId:  "customer-id",
			req:         dto.CreateCustomerSiteRequest{Name: "Cabang Bekasi"},
			existing:    []model.CustomerSite{existingSite},
			wantCreated: 1,
		},
		{
			name:        "explicit primary clears the previous primary",
			customerId:  "customer-id",
			req:         dto.CreateCustomerSiteRequest{Name: "Cabang Bekasi", IsPrimary: true},
			existing:    []model.CustomerSite{existingSite},
			wantPrimary: true,
			wantCleared: true,
			wantCreated: 1,
		},
		{
			name:        "rejects unknown customer",
			customerId:  "missing-id",
			req:         dto.CreateCustomerSiteRequest{Name: "Cabang"},
			wantCode:    404,
			wantCreated: 0,
		},
		{
			name:        "rejects blank name",
			customerId:  "customer-id",
			req:         dto.CreateCustomerSiteRequest{Name: "   "},
			wantCode:    400,
			wantCreated: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			customerRepo, siteRepo := newCustomerSiteTestRepos(c.existing...)
			usecase := newFakeCustomerSiteUsecase(customerRepo, siteRepo)

			err := usecase.Create("actor-id", c.customerId, &c.req)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Create() error = %v, want status %d", err, c.wantCode)
				}
				if len(siteRepo.created) != c.wantCreated {
					t.Fatalf("Create() persisted %d sites, want %d", len(siteRepo.created), c.wantCreated)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() error = %v, want nil", err.Error())
			}
			if len(siteRepo.created) != 1 {
				t.Fatalf("Create() persisted %d sites, want 1", len(siteRepo.created))
			}
			site := siteRepo.created[0]
			if site.CustomerId != c.customerId {
				t.Fatalf("Create() customer_id = %q, want %q", site.CustomerId, c.customerId)
			}
			if site.IsPrimary != c.wantPrimary {
				t.Fatalf("Create() is_primary = %v, want %v", site.IsPrimary, c.wantPrimary)
			}
			if !site.IsActive {
				t.Fatalf("Create() is_active = false, want true")
			}
			if site.Name != strings.TrimSpace(c.req.Name) {
				t.Fatalf("Create() name = %q, want %q", site.Name, strings.TrimSpace(c.req.Name))
			}
			if (siteRepo.clearedFor != "") != c.wantCleared {
				t.Fatalf("Create() clear-primary called = %v, want %v", siteRepo.clearedFor != "", c.wantCleared)
			}
		})
	}
}

func TestCustomerSiteUsecaseUpdateGuards(t *testing.T) {
	primary := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-1"}, CustomerId: "customer-id", Name: "Bandung", IsPrimary: true, IsActive: true}
	otherCustomer := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-x"}, CustomerId: "other-id", Name: "Surabaya", IsActive: true}

	cases := []struct {
		name        string
		id          string
		req         dto.UpdateCustomerSiteRequest
		wantCode    int
		wantPrimary bool
		wantCleared bool
	}{
		{
			name:        "promotes a site to primary",
			id:          "site-2",
			req:         dto.UpdateCustomerSiteRequest{Name: "Bekasi", IsPrimary: true, IsActive: true},
			wantPrimary: true,
			wantCleared: true,
		},
		{
			name:     "rejects site of another customer",
			id:       "site-x",
			req:      dto.UpdateCustomerSiteRequest{Name: "Surabaya", IsActive: true},
			wantCode: 404,
		},
		{
			name:     "rejects deactivating the primary site",
			id:       "site-1",
			req:      dto.UpdateCustomerSiteRequest{Name: "Bandung", IsPrimary: true, IsActive: false},
			wantCode: 400,
		},
		{
			name:     "rejects blank name",
			id:       "site-2",
			req:      dto.UpdateCustomerSiteRequest{Name: "  ", IsActive: true},
			wantCode: 400,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			second := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-2"}, CustomerId: "customer-id", Name: "Bekasi", IsActive: true}
			customerRepo, siteRepo := newCustomerSiteTestRepos(primary, second)
			siteRepo.byId["site-x"] = &otherCustomer
			usecase := newFakeCustomerSiteUsecase(customerRepo, siteRepo)

			err := usecase.Update("actor-id", "customer-id", c.id, &c.req)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Update() error = %v, want status %d", err, c.wantCode)
				}
				if siteRepo.updated != nil {
					t.Fatalf("Update() persisted a site, want none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() error = %v, want nil", err.Error())
			}
			if siteRepo.updated == nil {
				t.Fatalf("Update() persisted nothing")
			}
			if siteRepo.updated.IsPrimary != c.wantPrimary {
				t.Fatalf("Update() is_primary = %v, want %v", siteRepo.updated.IsPrimary, c.wantPrimary)
			}
			if !siteRepo.updated.IsActive {
				t.Fatalf("Update() is_active = false, want true")
			}
			if (siteRepo.clearedFor != "") != c.wantCleared {
				t.Fatalf("Update() clear-primary called = %v, want %v", siteRepo.clearedFor != "", c.wantCleared)
			}
		})
	}
}

func TestCustomerSiteUsecaseDeleteGuards(t *testing.T) {
	own := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-1"}, CustomerId: "customer-id", Name: "Bandung", IsPrimary: true, IsActive: true}
	foreign := model.CustomerSite{BaseModel: model.BaseModel{Id: "site-x"}, CustomerId: "other-id", Name: "Surabaya", IsActive: true}

	cases := []struct {
		name        string
		id          string
		wantCode    int
		wantDeleted string
	}{
		{"deletes own site", "site-1", 0, "site-1"},
		{"rejects site of another customer", "site-x", 404, ""},
		{"rejects unknown site", "missing-id", 404, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			customerRepo, siteRepo := newCustomerSiteTestRepos(own)
			siteRepo.byId["site-x"] = &foreign
			usecase := newFakeCustomerSiteUsecase(customerRepo, siteRepo)

			err := usecase.Delete("actor-id", "customer-id", c.id)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Delete() error = %v, want status %d", err, c.wantCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("Delete() error = %v, want nil", err.Error())
			}
			if siteRepo.deletedId != c.wantDeleted {
				t.Fatalf("Delete() deleted %q, want %q", siteRepo.deletedId, c.wantDeleted)
			}
		})
	}
}

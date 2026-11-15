package usecase

import (
	"testing"

	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
)

type fakeCustomerRepo struct {
	byId    map[string]*model.Customer
	created []model.Customer
	updated *model.Customer
}

func (f *fakeCustomerRepo) FindAll(page, limit int, search string) ([]model.Customer, int64, global.ErrorResponse) {
	return nil, 0, nil
}

func (f *fakeCustomerRepo) FindById(id string) (*model.Customer, global.ErrorResponse) {
	if customer, ok := f.byId[id]; ok {
		return customer, nil
	}
	return nil, global.NotFoundError("Customer not found")
}

func (f *fakeCustomerRepo) FindByIdForUpdate(tx helper.Tx, id string) (*model.Customer, global.ErrorResponse) {
	return f.FindById(id)
}

func (f *fakeCustomerRepo) FindByCode(code string) (*model.Customer, global.ErrorResponse) {
	return nil, global.NotFoundError("Customer not found")
}

func (f *fakeCustomerRepo) Create(tx helper.Tx, customer *model.Customer) global.ErrorResponse {
	f.created = append(f.created, *customer)
	return nil
}

func (f *fakeCustomerRepo) Update(tx helper.Tx, customer *model.Customer) global.ErrorResponse {
	f.updated = customer
	return nil
}

func (f *fakeCustomerRepo) AdjustOutstanding(tx helper.Tx, customerId string, delta int) global.ErrorResponse {
	return nil
}

func (f *fakeCustomerRepo) AdjustOutstandingBalance(tx helper.Tx, customerId string, delta float64) global.ErrorResponse {
	return nil
}

func newFakeCustomerUsecase(customerRepo *fakeCustomerRepo, topRepo *fakeTermOfPaymentRepo) CustomerUsecase {
	return NewCustomerUsecase(fakeTxManager{}, customerRepo, topRepo, fakeAuditLogRepo{})
}

func newCustomerTestTopRepo() *fakeTermOfPaymentRepo {
	active := &model.TermOfPayment{BaseModel: model.BaseModel{Id: "top-30"}, Code: "NET30", Name: "Net 30", Days: 30, IsActive: true}
	inactive := &model.TermOfPayment{BaseModel: model.BaseModel{Id: "top-off"}, Code: "OLD", Name: "Legacy", Days: 10, IsActive: false}
	return &fakeTermOfPaymentRepo{byId: map[string]*model.TermOfPayment{"top-30": active, "top-off": inactive}}
}

func TestCustomerUsecaseCreateProfileFields(t *testing.T) {
	cases := []struct {
		name      string
		topId     string
		wantCode  int
		wantTopId string
	}{
		{"persists new profile fields without TOP", "", 0, ""},
		{"stores active term of payment", "top-30", 0, "top-30"},
		{"ignores surrounding spaces on TOP id", " top-30 ", 0, "top-30"},
		{"rejects unknown term of payment", "missing-id", 404, ""},
		{"rejects inactive term of payment", "top-off", 400, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			customerRepo := &fakeCustomerRepo{}
			usecase := newFakeCustomerUsecase(customerRepo, newCustomerTestTopRepo())

			req := dto.CreateCustomerRequest{
				Name:            "PT Maju Industri",
				Npwp:            "01.234.567.8-901.000",
				PostalCode:      "40115",
				CreditLimit:     25000000,
				TermOfPaymentId: c.topId,
				SalesPicName:    "Budi",
				SalesPicPhone:   "081200000000",
				SalesPicEmail:   "budi@maju.co.id",
			}

			err := usecase.Create("actor-id", &req)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Create() error = %v, want status %d", err, c.wantCode)
				}
				if len(customerRepo.created) != 0 {
					t.Fatalf("Create() persisted %d customers, want 0", len(customerRepo.created))
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() error = %v, want nil", err.Error())
			}
			if len(customerRepo.created) != 1 {
				t.Fatalf("Create() persisted %d customers, want 1", len(customerRepo.created))
			}

			got := customerRepo.created[0]
			if got.Npwp != req.Npwp || got.PostalCode != req.PostalCode || got.CreditLimit != req.CreditLimit {
				t.Fatalf("Create() persisted %+v, want npwp/postal_code/credit_limit from request", got)
			}
			if got.SalesPicName != req.SalesPicName || got.SalesPicPhone != req.SalesPicPhone || got.SalesPicEmail != req.SalesPicEmail {
				t.Fatalf("Create() persisted sales PIC %+v, want request values", got)
			}
			if got.PaymentTermDays != 30 {
				t.Fatalf("Create() payment_term_days = %d, want default 30", got.PaymentTermDays)
			}
			switch {
			case c.wantTopId == "" && got.TermOfPaymentId != nil:
				t.Fatalf("Create() term_of_payment_id = %q, want nil", *got.TermOfPaymentId)
			case c.wantTopId != "" && (got.TermOfPaymentId == nil || *got.TermOfPaymentId != c.wantTopId):
				t.Fatalf("Create() term_of_payment_id = %v, want %q", got.TermOfPaymentId, c.wantTopId)
			}
		})
	}
}

func TestCustomerUsecaseUpdateTermOfPayment(t *testing.T) {
	cases := []struct {
		name      string
		topId     string
		wantCode  int
		wantTopId string
	}{
		{"sets term of payment", "top-30", 0, "top-30"},
		{"clears term of payment with empty id", "", 0, ""},
		{"rejects unknown term of payment", "missing-id", 404, ""},
		{"rejects inactive term of payment", "top-off", 400, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			existing := &model.Customer{BaseModel: model.BaseModel{Id: "customer-id"}, Name: "PT Lama", PaymentTermDays: 45}
			customerRepo := &fakeCustomerRepo{byId: map[string]*model.Customer{"customer-id": existing}}
			usecase := newFakeCustomerUsecase(customerRepo, newCustomerTestTopRepo())

			req := dto.UpdateCustomerRequest{Name: "PT Baru", CreditLimit: 5000000, PostalCode: "40123", TermOfPaymentId: c.topId}
			err := usecase.Update("actor-id", "customer-id", &req)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Update() error = %v, want status %d", err, c.wantCode)
				}
				if customerRepo.updated != nil {
					t.Fatalf("Update() persisted a customer, want none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() error = %v, want nil", err.Error())
			}
			if customerRepo.updated == nil {
				t.Fatalf("Update() persisted nothing")
			}
			got := customerRepo.updated
			if c.wantTopId == "" && got.TermOfPaymentId != nil {
				t.Fatalf("Update() term_of_payment_id = %q, want nil", *got.TermOfPaymentId)
			}
			if c.wantTopId != "" && (got.TermOfPaymentId == nil || *got.TermOfPaymentId != c.wantTopId) {
				t.Fatalf("Update() term_of_payment_id = %v, want %q", got.TermOfPaymentId, c.wantTopId)
			}
			// TOP is not wired to invoicing yet, so the existing grace period must not move.
			if got.PaymentTermDays != 45 {
				t.Fatalf("Update() payment_term_days = %d, want unchanged 45", got.PaymentTermDays)
			}
			if got.CreditLimit != req.CreditLimit || got.PostalCode != req.PostalCode {
				t.Fatalf("Update() persisted %+v, want credit_limit/postal_code from request", got)
			}
		})
	}
}

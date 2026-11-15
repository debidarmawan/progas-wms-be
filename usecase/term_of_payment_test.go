package usecase

import (
	"testing"

	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
)

type fakeTermOfPaymentRepo struct {
	byId      map[string]*model.TermOfPayment
	byCode    map[string]*model.TermOfPayment
	created   []model.TermOfPayment
	updated   *model.TermOfPayment
	deletedId string
}

func (f *fakeTermOfPaymentRepo) FindAll(page, limit int, search string) ([]model.TermOfPayment, int64, global.ErrorResponse) {
	return nil, 0, nil
}

func (f *fakeTermOfPaymentRepo) FindById(id string) (*model.TermOfPayment, global.ErrorResponse) {
	if top, ok := f.byId[id]; ok {
		return top, nil
	}
	return nil, global.NotFoundError("Term of payment not found")
}

func (f *fakeTermOfPaymentRepo) FindByCode(code string) (*model.TermOfPayment, global.ErrorResponse) {
	if top, ok := f.byCode[code]; ok {
		return top, nil
	}
	return nil, global.NotFoundError("Term of payment not found")
}

func (f *fakeTermOfPaymentRepo) Create(tx helper.Tx, top *model.TermOfPayment) global.ErrorResponse {
	top.Id = "created-id"
	f.created = append(f.created, *top)
	return nil
}

func (f *fakeTermOfPaymentRepo) Update(tx helper.Tx, top *model.TermOfPayment) global.ErrorResponse {
	f.updated = top
	return nil
}

func (f *fakeTermOfPaymentRepo) Delete(tx helper.Tx, id string) global.ErrorResponse {
	f.deletedId = id
	return nil
}

func newFakeTermOfPaymentUsecase(repo *fakeTermOfPaymentRepo) TermOfPaymentUsecase {
	return NewTermOfPaymentUsecase(fakeTxManager{}, repo, fakeAuditLogRepo{})
}

func TestTermOfPaymentUsecaseCreate(t *testing.T) {
	existing := &model.TermOfPayment{BaseModel: model.BaseModel{Id: "existing-id"}, Code: "NET30", Name: "Net 30", Days: 30}

	cases := []struct {
		name     string
		req      dto.CreateTermOfPaymentRequest
		wantCode int
		want     *model.TermOfPayment
	}{
		{
			name: "normalizes code and trims name, active by default",
			req:  dto.CreateTermOfPaymentRequest{Code: " net14 ", Name: " Net 14 Hari ", Days: 14},
			want: &model.TermOfPayment{Code: "NET14", Name: "Net 14 Hari", Days: 14, IsActive: true},
		},
		{
			name: "allows cash term with zero days",
			req:  dto.CreateTermOfPaymentRequest{Code: "CASH", Name: "Cash", Days: 0},
			want: &model.TermOfPayment{Code: "CASH", Name: "Cash", Days: 0, IsActive: true},
		},
		{
			name:     "rejects duplicate code",
			req:      dto.CreateTermOfPaymentRequest{Code: "net30", Name: "Net 30", Days: 30},
			wantCode: 400,
		},
		{
			name:     "rejects blank code",
			req:      dto.CreateTermOfPaymentRequest{Code: "  ", Name: "Net 30", Days: 30},
			wantCode: 400,
		},
		{
			name:     "rejects blank name",
			req:      dto.CreateTermOfPaymentRequest{Code: "NET7", Name: "  ", Days: 7},
			wantCode: 400,
		},
		{
			name:     "rejects negative days",
			req:      dto.CreateTermOfPaymentRequest{Code: "NEG", Name: "Negative", Days: -1},
			wantCode: 400,
		},
		{
			name:     "rejects days beyond one year",
			req:      dto.CreateTermOfPaymentRequest{Code: "LONG", Name: "Too long", Days: 366},
			wantCode: 400,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := &fakeTermOfPaymentRepo{
				byCode: map[string]*model.TermOfPayment{"NET30": existing},
				byId:   map[string]*model.TermOfPayment{"existing-id": existing},
			}
			usecase := newFakeTermOfPaymentUsecase(repo)

			err := usecase.Create("actor-id", &c.req)

			if c.wantCode != 0 {
				if err == nil {
					t.Fatalf("Create() error = nil, want status %d", c.wantCode)
				}
				if err.GetCode() != c.wantCode {
					t.Fatalf("Create() status = %d, want %d", err.GetCode(), c.wantCode)
				}
				if len(repo.created) != 0 {
					t.Fatalf("Create() persisted %d terms, want 0", len(repo.created))
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() error = %v, want nil", err.Error())
			}
			if len(repo.created) != 1 {
				t.Fatalf("Create() persisted %d terms, want 1", len(repo.created))
			}
			got := repo.created[0]
			if got.Code != c.want.Code || got.Name != c.want.Name || got.Days != c.want.Days || got.IsActive != c.want.IsActive {
				t.Fatalf("Create() persisted %+v, want %+v", got, *c.want)
			}
		})
	}
}

func TestTermOfPaymentUsecaseUpdate(t *testing.T) {
	cases := []struct {
		name        string
		days        int
		wantCode    int
		wantDays    int
		wantUpdated bool
	}{
		{"updates to cash term", 0, 0, 0, true},
		{"updates to 365 days", 365, 0, 365, true},
		{"rejects days beyond one year", 366, 400, 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			top := &model.TermOfPayment{BaseModel: model.BaseModel{Id: "top-id"}, Code: "NET30", Name: "Net 30", Days: 30, IsActive: true}
			repo := &fakeTermOfPaymentRepo{
				byId:   map[string]*model.TermOfPayment{"top-id": top},
				byCode: map[string]*model.TermOfPayment{"NET30": top},
			}
			usecase := newFakeTermOfPaymentUsecase(repo)

			req := dto.UpdateTermOfPaymentRequest{Code: "NET30", Name: "Net 30", Days: c.days}
			err := usecase.Update("actor-id", "top-id", &req)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Update() error = %v, want status %d", err, c.wantCode)
				}
			} else if err != nil {
				t.Fatalf("Update() error = %v, want nil", err.Error())
			}

			updated := repo.updated != nil
			if updated != c.wantUpdated {
				t.Fatalf("Update() persisted = %v, want %v", updated, c.wantUpdated)
			}
			if updated && repo.updated.Days != c.wantDays {
				t.Fatalf("Update() days = %d, want %d", repo.updated.Days, c.wantDays)
			}
		})
	}
}

func TestTermOfPaymentUsecaseDelete(t *testing.T) {
	cases := []struct {
		name        string
		id          string
		wantCode    int
		wantDeleted bool
	}{
		{"deletes existing term", "top-id", 0, true},
		{"rejects unknown term", "missing-id", 404, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			top := &model.TermOfPayment{BaseModel: model.BaseModel{Id: "top-id"}, Code: "CASH", Name: "Cash"}
			repo := &fakeTermOfPaymentRepo{byId: map[string]*model.TermOfPayment{"top-id": top}}
			usecase := newFakeTermOfPaymentUsecase(repo)

			err := usecase.Delete("actor-id", c.id)

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Delete() error = %v, want status %d", err, c.wantCode)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v, want nil", err.Error())
			}

			deleted := repo.deletedId == "top-id"
			if deleted != c.wantDeleted {
				t.Fatalf("Delete() deleted = %v, want %v", deleted, c.wantDeleted)
			}
		})
	}
}

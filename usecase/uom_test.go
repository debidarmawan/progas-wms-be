package usecase

import (
	"testing"

	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type fakeTx struct{}

func (fakeTx) Rollback()     {}
func (fakeTx) Commit() error { return nil }
func (fakeTx) CheckPanic()   {}
func (fakeTx) Get() *gorm.DB { return nil }

type fakeTxManager struct{}

func (fakeTxManager) New() helper.Tx { return fakeTx{} }

type fakeUomRepo struct {
	byId      map[string]*model.Uom
	byCode    map[string]*model.Uom
	created   []model.Uom
	deletedId string
}

func (f *fakeUomRepo) FindAll(page, limit int, search string) ([]model.Uom, int64, global.ErrorResponse) {
	return nil, 0, nil
}

func (f *fakeUomRepo) FindById(id string) (*model.Uom, global.ErrorResponse) {
	if uom, ok := f.byId[id]; ok {
		return uom, nil
	}
	return nil, global.NotFoundError("UOM not found")
}

func (f *fakeUomRepo) FindByCode(code string) (*model.Uom, global.ErrorResponse) {
	if uom, ok := f.byCode[code]; ok {
		return uom, nil
	}
	return nil, global.NotFoundError("UOM not found")
}

func (f *fakeUomRepo) Create(tx helper.Tx, uom *model.Uom) global.ErrorResponse {
	uom.Id = "created-id"
	f.created = append(f.created, *uom)
	return nil
}

func (f *fakeUomRepo) Update(tx helper.Tx, uom *model.Uom) global.ErrorResponse {
	return nil
}

func (f *fakeUomRepo) Delete(tx helper.Tx, id string) global.ErrorResponse {
	f.deletedId = id
	return nil
}

type fakeMasterItemRepo struct {
	usedByUomId int64
}

func (f *fakeMasterItemRepo) FindAll(page, limit int, search, sortBy, sortOrder, itemType, gasType string, isSerialized *bool) ([]model.MasterItem, int64, global.ErrorResponse) {
	return nil, 0, nil
}

func (f *fakeMasterItemRepo) FindById(id string) (*model.MasterItem, global.ErrorResponse) {
	return nil, global.NotFoundError("Master item not found")
}

func (f *fakeMasterItemRepo) FindBySKU(sku string) (*model.MasterItem, global.ErrorResponse) {
	return nil, global.NotFoundError("Master item not found")
}

func (f *fakeMasterItemRepo) CountByUomId(uomId string) (int64, global.ErrorResponse) {
	return f.usedByUomId, nil
}

func (f *fakeMasterItemRepo) Create(tx helper.Tx, item *model.MasterItem) global.ErrorResponse {
	return nil
}

func (f *fakeMasterItemRepo) Update(tx helper.Tx, item *model.MasterItem) global.ErrorResponse {
	return nil
}

type fakeAuditLogRepo struct{}

func (fakeAuditLogRepo) Create(log *model.AuditLog) global.ErrorResponse { return nil }
func (fakeAuditLogRepo) Log(userId, action, objectType, objectId string, details any) global.ErrorResponse {
	return nil
}
func (fakeAuditLogRepo) FindByObject(objectType, objectId string) ([]model.AuditLog, global.ErrorResponse) {
	return nil, nil
}

func newFakeUomUsecase(uomRepo *fakeUomRepo, masterItemRepo *fakeMasterItemRepo) UomUsecase {
	return NewUomUsecase(fakeTxManager{}, uomRepo, masterItemRepo, fakeAuditLogRepo{})
}

func TestUomUsecaseCreate(t *testing.T) {
	existing := &model.Uom{BaseModel: model.BaseModel{Id: "existing-id"}, Code: "KG", Name: "Kilogram"}

	cases := []struct {
		name     string
		req      dto.CreateUomRequest
		wantCode int
		wantUom  *model.Uom
	}{
		{
			name:    "normalizes code to uppercase and trims",
			req:     dto.CreateUomRequest{Code: " pcs ", Name: " Pieces "},
			wantUom: &model.Uom{Code: "PCS", Name: "Pieces"},
		},
		{
			name:     "rejects duplicate code",
			req:      dto.CreateUomRequest{Code: "kg", Name: "Kilogram"},
			wantCode: 400,
		},
		{
			name:     "rejects blank code",
			req:      dto.CreateUomRequest{Code: "  ", Name: "Kilogram"},
			wantCode: 400,
		},
		{
			name:     "rejects blank name",
			req:      dto.CreateUomRequest{Code: "M3", Name: "  "},
			wantCode: 400,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			uomRepo := &fakeUomRepo{
				byCode: map[string]*model.Uom{"KG": existing},
				byId:   map[string]*model.Uom{"existing-id": existing},
			}
			usecase := newFakeUomUsecase(uomRepo, &fakeMasterItemRepo{})

			err := usecase.Create("actor-id", &c.req)

			if c.wantCode != 0 {
				if err == nil {
					t.Fatalf("Create() error = nil, want status %d", c.wantCode)
				}
				if err.GetCode() != c.wantCode {
					t.Fatalf("Create() status = %d, want %d", err.GetCode(), c.wantCode)
				}
				if len(uomRepo.created) != 0 {
					t.Fatalf("Create() persisted %d uoms, want 0", len(uomRepo.created))
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() error = %v, want nil", err.Error())
			}
			if len(uomRepo.created) != 1 {
				t.Fatalf("Create() persisted %d uoms, want 1", len(uomRepo.created))
			}
			if got := uomRepo.created[0]; got.Code != c.wantUom.Code || got.Name != c.wantUom.Name {
				t.Fatalf("Create() persisted %+v, want code %q name %q", got, c.wantUom.Code, c.wantUom.Name)
			}
		})
	}
}

func TestUomUsecaseDelete(t *testing.T) {
	cases := []struct {
		name        string
		usedByUomId int64
		wantCode    int
		wantDeleted bool
	}{
		{"deletes unused uom", 0, 0, true},
		{"rejects uom still referenced by master items", 2, 400, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			uom := &model.Uom{BaseModel: model.BaseModel{Id: "uom-id"}, Code: "PCS", Name: "Pieces"}
			uomRepo := &fakeUomRepo{byId: map[string]*model.Uom{"uom-id": uom}}
			usecase := newFakeUomUsecase(uomRepo, &fakeMasterItemRepo{usedByUomId: c.usedByUomId})

			err := usecase.Delete("actor-id", "uom-id")

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Delete() error = %v, want status %d", err, c.wantCode)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v, want nil", err.Error())
			}

			deleted := uomRepo.deletedId == "uom-id"
			if deleted != c.wantDeleted {
				t.Fatalf("Delete() deleted = %v, want %v", deleted, c.wantDeleted)
			}
		})
	}
}

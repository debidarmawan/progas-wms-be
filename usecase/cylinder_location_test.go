package usecase

import (
	"testing"

	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
)

type fakeCylinderLocationRepo struct {
	byId      map[string]*model.CylinderLocation
	byCode    map[string]*model.CylinderLocation
	used      int64
	created   []model.CylinderLocation
	updated   *model.CylinderLocation
	deletedId string
}

func (f *fakeCylinderLocationRepo) FindAll(page, limit int, search string) ([]model.CylinderLocation, int64, global.ErrorResponse) {
	return nil, 0, nil
}

func (f *fakeCylinderLocationRepo) FindById(id string) (*model.CylinderLocation, global.ErrorResponse) {
	if location, ok := f.byId[id]; ok {
		return location, nil
	}
	return nil, global.NotFoundError("Cylinder location not found")
}

func (f *fakeCylinderLocationRepo) FindByCode(code string) (*model.CylinderLocation, global.ErrorResponse) {
	if location, ok := f.byCode[code]; ok {
		return location, nil
	}
	return nil, global.NotFoundError("Cylinder location not found")
}

func (f *fakeCylinderLocationRepo) CountCylindersByLocationId(id string) (int64, global.ErrorResponse) {
	return f.used, nil
}

func (f *fakeCylinderLocationRepo) Create(tx helper.Tx, location *model.CylinderLocation) global.ErrorResponse {
	location.Id = "created-id"
	f.created = append(f.created, *location)
	return nil
}

func (f *fakeCylinderLocationRepo) Update(tx helper.Tx, location *model.CylinderLocation) global.ErrorResponse {
	f.updated = location
	return nil
}

func (f *fakeCylinderLocationRepo) Delete(tx helper.Tx, id string) global.ErrorResponse {
	f.deletedId = id
	return nil
}

func newFakeCylinderLocationUsecase(repo *fakeCylinderLocationRepo) CylinderLocationUsecase {
	return NewCylinderLocationUsecase(fakeTxManager{}, repo, fakeAuditLogRepo{})
}

func TestCylinderLocationUsecaseCreate(t *testing.T) {
	existing := &model.CylinderLocation{BaseModel: model.BaseModel{Id: "existing-id"}, Code: "WH-1", Name: "Gudang Utama"}

	cases := []struct {
		name       string
		req        dto.CreateCylinderLocationRequest
		wantCode   int
		wantCodeX  string
		wantName   string
		wantType   string
		wantActive bool
	}{
		{
			name:       "normalizes code and name",
			req:        dto.CreateCylinderLocationRequest{Code: " wh-2 ", Name: " Gudang Cabang ", Type: "warehouse"},
			wantCodeX:  "WH-2",
			wantName:   "Gudang Cabang",
			wantType:   "WAREHOUSE",
			wantActive: true,
		},
		{
			name:     "rejects duplicate code",
			req:      dto.CreateCylinderLocationRequest{Code: "wh-1", Name: "Gudang Lain", Type: "WAREHOUSE"},
			wantCode: 400,
		},
		{
			name:     "rejects non-internal type",
			req:      dto.CreateCylinderLocationRequest{Code: "CUST-1", Name: "Customer", Type: "CUSTOMER"},
			wantCode: 400,
		},
		{
			name:     "rejects blank name",
			req:      dto.CreateCylinderLocationRequest{Code: "FS-1", Name: "  ", Type: "FILLING_STATION"},
			wantCode: 400,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := &fakeCylinderLocationRepo{
				byId:   map[string]*model.CylinderLocation{"existing-id": existing},
				byCode: map[string]*model.CylinderLocation{"WH-1": existing},
			}
			usecase := newFakeCylinderLocationUsecase(repo)

			err := usecase.Create("actor-id", &c.req)

			if c.wantCode != 0 {
				if err == nil {
					t.Fatalf("Create() error = nil, want status %d", c.wantCode)
				}
				if err.GetCode() != c.wantCode {
					t.Fatalf("Create() status = %d, want %d", err.GetCode(), c.wantCode)
				}
				if len(repo.created) != 0 {
					t.Fatalf("Create() persisted %d locations, want 0", len(repo.created))
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() error = %v, want nil", err.Error())
			}
			if len(repo.created) != 1 {
				t.Fatalf("Create() persisted %d locations, want 1", len(repo.created))
			}
			got := repo.created[0]
			if got.Code != c.wantCodeX || got.Name != c.wantName || string(got.Type) != c.wantType || got.IsActive != c.wantActive {
				t.Fatalf("Create() persisted %+v, want code %q name %q type %q active %v", got, c.wantCodeX, c.wantName, c.wantType, c.wantActive)
			}
		})
	}
}

func TestCylinderLocationUsecaseUpdateIsActive(t *testing.T) {
	cases := []struct {
		name     string
		isActive *bool
		want     bool
	}{
		{"keeps active when omitted", nil, true},
		{"deactivates on request", boolPtr(false), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			location := &model.CylinderLocation{BaseModel: model.BaseModel{Id: "loc-1"}, Code: "WH-1", Name: "Gudang", Type: "WAREHOUSE", IsActive: true}
			repo := &fakeCylinderLocationRepo{
				byId:   map[string]*model.CylinderLocation{"loc-1": location},
				byCode: map[string]*model.CylinderLocation{"WH-1": location},
			}
			usecase := newFakeCylinderLocationUsecase(repo)

			err := usecase.Update("actor-id", "loc-1", &dto.UpdateCylinderLocationRequest{
				Code:     "WH-1",
				Name:     "Gudang",
				Type:     "WAREHOUSE",
				IsActive: c.isActive,
			})
			if err != nil {
				t.Fatalf("Update() error = %v, want nil", err.Error())
			}
			if repo.updated == nil {
				t.Fatalf("Update() persisted nothing")
			}
			if repo.updated.IsActive != c.want {
				t.Fatalf("Update() is_active = %v, want %v", repo.updated.IsActive, c.want)
			}
		})
	}
}

func TestCylinderLocationUsecaseDeleteGuards(t *testing.T) {
	cases := []struct {
		name        string
		used        int64
		wantCode    int
		wantDeleted bool
	}{
		{"deletes unused location", 0, 0, true},
		{"rejects location still assigned to cylinders", 3, 400, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			location := &model.CylinderLocation{BaseModel: model.BaseModel{Id: "loc-1"}, Code: "WH-1", Name: "Gudang", Type: "WAREHOUSE", IsActive: true}
			repo := &fakeCylinderLocationRepo{byId: map[string]*model.CylinderLocation{"loc-1": location}, used: c.used}
			usecase := newFakeCylinderLocationUsecase(repo)

			err := usecase.Delete("actor-id", "loc-1")

			if c.wantCode != 0 {
				if err == nil || err.GetCode() != c.wantCode {
					t.Fatalf("Delete() error = %v, want status %d", err, c.wantCode)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v, want nil", err.Error())
			}

			deleted := repo.deletedId == "loc-1"
			if deleted != c.wantDeleted {
				t.Fatalf("Delete() deleted = %v, want %v", deleted, c.wantDeleted)
			}
		})
	}
}

func boolPtr(value bool) *bool {
	return &value
}

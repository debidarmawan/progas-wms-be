package usecase

import (
	"testing"

	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func TestOptionalId(t *testing.T) {
	const accountId = "account-id"
	cases := []struct {
		name string
		in   string
		want *string
	}{
		{"blank becomes nil", "", nil},
		{"keeps provided id", accountId, stringPtr(accountId)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := optionalId(c.in)
			if c.want == nil {
				if got != nil {
					t.Fatalf("optionalId(%q) = %q, want nil", c.in, *got)
				}
				return
			}
			if got == nil || *got != *c.want {
				t.Fatalf("optionalId(%q) = %v, want %q", c.in, got, *c.want)
			}
		})
	}
}

// Regression: a blank COA mapping was saved as "" and violated fk_master_item_inventory_account.
func TestMasterItemUsecaseUpdateBlankAccountMappingStoresNull(t *testing.T) {
	item := &model.MasterItem{
		BaseModel:    model.BaseModel{Id: "item-id"},
		Name:         "Oxygen",
		SKU:          "GAS-O2",
		ItemType:     "gas",
		GasType:      "O2",
		IsSerialized: true,
	}
	repo := &fakeMasterItemRepo{byId: map[string]*model.MasterItem{"item-id": item}}
	usecase := NewMasterItemUsecase(fakeTxManager{}, repo, nil, &fakeUomRepo{}, nil, fakeAuditLogRepo{})

	err := usecase.Update("actor-id", "item-id", &dto.UpdateMasterItemRequest{Name: "Oxygen", ItemType: "gas", GasType: "O2"})
	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err.Error())
	}
	if repo.updated == nil {
		t.Fatal("Update() persisted no item")
	}
	if repo.updated.InventoryAccountId != nil {
		t.Fatalf("InventoryAccountId = %q, want nil (NULL)", *repo.updated.InventoryAccountId)
	}
	if repo.updated.RevenueAccountId != nil || repo.updated.CogsAccountId != nil {
		t.Fatalf("RevenueAccountId/CogsAccountId = %v/%v, want nil", repo.updated.RevenueAccountId, repo.updated.CogsAccountId)
	}
}

func stringPtr(value string) *string {
	return &value
}

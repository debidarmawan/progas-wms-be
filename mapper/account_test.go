package mapper

import (
	"testing"

	"progas-wms-be/enum"
	"progas-wms-be/model"
)

func TestToAccountResponseParent(t *testing.T) {
	parentId := "parent-id"
	cases := []struct {
		name           string
		account        model.Account
		wantParentId   string
		wantParentCode string
	}{
		{
			name:    "root account has no parent id",
			account: model.Account{Code: "1-1500", Type: enum.AccountTypeAsset},
		},
		{
			name:         "child account keeps parent id",
			account:      model.Account{Code: "1-1501", Type: enum.AccountTypeAsset, ParentId: &parentId},
			wantParentId: parentId,
		},
		{
			name: "preloaded parent exposes parent code",
			account: model.Account{
				Code:     "1-1501",
				Type:     enum.AccountTypeAsset,
				ParentId: &parentId,
				Parent:   &model.Account{Code: "1-1500"},
			},
			wantParentId:   parentId,
			wantParentCode: "1-1500",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ToAccountResponse(&c.account)
			if got.ParentId != c.wantParentId {
				t.Errorf("ParentId = %q, want %q", got.ParentId, c.wantParentId)
			}
			if got.ParentCode != c.wantParentCode {
				t.Errorf("ParentCode = %q, want %q", got.ParentCode, c.wantParentCode)
			}
		})
	}
}

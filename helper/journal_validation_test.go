package helper

import (
	"progas-wms-be/dto"
	"testing"
)

func TestCreateJournalEntryRequestValidation(t *testing.T) {
	line := func(accountId string, debit, credit float64) dto.CreateJournalLineRequest {
		return dto.CreateJournalLineRequest{AccountId: accountId, Debit: debit, Credit: credit}
	}

	cases := []struct {
		name    string
		req     dto.CreateJournalEntryRequest
		wantErr bool
	}{
		{
			name: "valid two line entry",
			req: dto.CreateJournalEntryRequest{
				EntryDate:   "2026-11-09",
				Description: "Setoran modal",
				Lines:       []dto.CreateJournalLineRequest{line("acc-cash", 1000, 0), line("acc-equity", 0, 1000)},
			},
		},
		{
			name: "missing entry date rejected",
			req: dto.CreateJournalEntryRequest{
				Lines: []dto.CreateJournalLineRequest{line("acc-cash", 1000, 0), line("acc-equity", 0, 1000)},
			},
			wantErr: true,
		},
		{
			name: "bad date format rejected",
			req: dto.CreateJournalEntryRequest{
				EntryDate: "09-11-2026",
				Lines:     []dto.CreateJournalLineRequest{line("acc-cash", 1000, 0), line("acc-equity", 0, 1000)},
			},
			wantErr: true,
		},
		{
			name: "single line rejected",
			req: dto.CreateJournalEntryRequest{
				EntryDate: "2026-11-09",
				Lines:     []dto.CreateJournalLineRequest{line("acc-cash", 1000, 0)},
			},
			wantErr: true,
		},
		{
			name: "line without account rejected",
			req: dto.CreateJournalEntryRequest{
				EntryDate: "2026-11-09",
				Lines:     []dto.CreateJournalLineRequest{line("", 1000, 0), line("acc-equity", 0, 1000)},
			},
			wantErr: true,
		},
		{
			name: "negative debit rejected",
			req: dto.CreateJournalEntryRequest{
				EntryDate: "2026-11-09",
				Lines:     []dto.CreateJournalLineRequest{line("acc-cash", -1000, 0), line("acc-equity", 0, 1000)},
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			errs := ValidateStruct(&c.req)
			if c.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation errors, got none")
			}
			if !c.wantErr && len(errs) != 0 {
				t.Fatalf("unexpected validation errors: %+v", errs)
			}
		})
	}
}

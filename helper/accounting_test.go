package helper

import (
	"progas-wms-be/constant"
	"progas-wms-be/enum"
	"strings"
	"testing"
)

func TestValidateJournalLines(t *testing.T) {
	line := func(accountId string, debit, credit float64) JournalLineInput {
		return JournalLineInput{AccountId: accountId, Debit: debit, Credit: credit}
	}

	cases := []struct {
		name          string
		lines         []JournalLineInput
		wantErr       bool
		wantDebit     float64
		wantCredit    float64
		errorContains string
	}{
		{
			name:       "balanced debit and credit",
			lines:      []JournalLineInput{line("acc-cash", 150000, 0), line("acc-revenue", 0, 150000)},
			wantDebit:  150000,
			wantCredit: 150000,
		},
		{
			name: "multi line balanced",
			lines: []JournalLineInput{
				line("acc-cash", 100000, 0),
				line("acc-ar", 50000, 0),
				line("acc-revenue", 0, 150000),
			},
			wantDebit:  150000,
			wantCredit: 150000,
		},
		{
			name:    "single line rejected",
			lines:   []JournalLineInput{line("acc-cash", 1000, 0)},
			wantErr: true,
		},
		{
			name:    "empty lines rejected",
			lines:   nil,
			wantErr: true,
		},
		{
			name:          "unbalanced rejected",
			lines:         []JournalLineInput{line("acc-cash", 1000, 0), line("acc-revenue", 0, 900)},
			wantErr:       true,
			errorContains: "not balanced",
		},
		{
			name:          "both sides on one line rejected",
			lines:         []JournalLineInput{line("acc-cash", 1000, 1000), line("acc-revenue", 0, 1000)},
			wantErr:       true,
			errorContains: "only one of debit or credit",
		},
		{
			name:          "zero amount line rejected",
			lines:         []JournalLineInput{line("acc-cash", 0, 0), line("acc-revenue", 0, 1000)},
			wantErr:       true,
			errorContains: "debit or credit is required",
		},
		{
			name:          "negative amount rejected",
			lines:         []JournalLineInput{line("acc-cash", -1000, 0), line("acc-revenue", 0, -1000)},
			wantErr:       true,
			errorContains: "cannot be negative",
		},
		{
			name:          "missing account rejected",
			lines:         []JournalLineInput{line("", 1000, 0), line("acc-revenue", 0, 1000)},
			wantErr:       true,
			errorContains: "account is required",
		},
		{
			name:       "float noise within tolerance accepted",
			lines:      []JournalLineInput{line("acc-cash", 100.004, 0), line("acc-revenue", 0, 100)},
			wantDebit:  100,
			wantCredit: 100,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			debit, credit, err := ValidateJournalLines(c.lines)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if c.errorContains != "" {
					if got := err.Error(); !strings.Contains(got, c.errorContains) {
						t.Fatalf("error %q does not contain %q", got, c.errorContains)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if debit != c.wantDebit || credit != c.wantCredit {
				t.Fatalf("totals = (%.2f, %.2f), want (%.2f, %.2f)", debit, credit, c.wantDebit, c.wantCredit)
			}
		})
	}
}

func TestRoundMoney(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{0, 0},
		{1000.005, 1000.01},
		{1000.004, 1000},
		{-12.345, -12.35},
	}
	for _, c := range cases {
		if got := RoundMoney(c.in); got != c.want {
			t.Errorf("RoundMoney(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestValidateItemAccountType(t *testing.T) {
	slots := []struct {
		field string
		want  enum.AccountType
	}{
		{constant.ItemAccountInventory, enum.AccountTypeAsset},
		{constant.ItemAccountRevenue, enum.AccountTypeRevenue},
		{constant.ItemAccountCogs, enum.AccountTypeExpense},
	}
	allTypes := []enum.AccountType{
		enum.AccountTypeAsset,
		enum.AccountTypeLiability,
		enum.AccountTypeEquity,
		enum.AccountTypeRevenue,
		enum.AccountTypeExpense,
	}

	for _, slot := range slots {
		for _, accountType := range allTypes {
			t.Run(slot.field+"/"+string(accountType), func(t *testing.T) {
				err := ValidateItemAccountType(slot.field, accountType)
				if accountType == slot.want && err != nil {
					t.Errorf("expected %s accepted for %s, got %v", accountType, slot.field, err)
				}
				if accountType != slot.want && err == nil {
					t.Errorf("expected %s rejected for %s", accountType, slot.field)
				}
			})
		}
	}

	t.Run("unknown field rejected", func(t *testing.T) {
		if err := ValidateItemAccountType("mystery_account_id", enum.AccountTypeAsset); err == nil {
			t.Error("expected error for unknown mapping field")
		}
	})
}

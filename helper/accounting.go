package helper

import (
	"errors"
	"fmt"
	"math"
	"progas-wms-be/enum"
	"time"

	"github.com/google/uuid"
)

// journalBalanceTolerance absorbs float noise on 2-decimal money amounts.
const journalBalanceTolerance = 0.005

// AccountNormalBalance returns DEBIT/CREDIT side where the account increases.
func AccountNormalBalance(accountType enum.AccountType) string {
	switch accountType {
	case enum.AccountTypeAsset, enum.AccountTypeExpense:
		return "DEBIT"
	default:
		return "CREDIT"
	}
}

// JournalLineInput is the minimum data needed to post one journal line.
type JournalLineInput struct {
	AccountId string
	Debit     float64
	Credit    float64
}

func GenerateJournalEntryNumber() string {
	id, _ := uuid.NewV7()
	return fmt.Sprintf("JE-%s-%s", time.Now().Format("20060102"), id.String()[:8])
}

func RoundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

// ValidateJournalLines enforces double-entry rules: at least two lines, every
// line carries exactly one side (debit XOR credit), and both totals match.
func ValidateJournalLines(lines []JournalLineInput) (totalDebit, totalCredit float64, err error) {
	if len(lines) < 2 {
		return 0, 0, errors.New("journal must have at least two lines")
	}

	for i, line := range lines {
		debit := RoundMoney(line.Debit)
		credit := RoundMoney(line.Credit)
		if line.AccountId == "" {
			return 0, 0, fmt.Errorf("line %d: account is required", i+1)
		}
		if debit < 0 || credit < 0 {
			return 0, 0, fmt.Errorf("line %d: amount cannot be negative", i+1)
		}
		if debit > 0 && credit > 0 {
			return 0, 0, fmt.Errorf("line %d: only one of debit or credit may be filled", i+1)
		}
		if debit == 0 && credit == 0 {
			return 0, 0, fmt.Errorf("line %d: debit or credit is required", i+1)
		}
		totalDebit += debit
		totalCredit += credit
	}

	totalDebit = RoundMoney(totalDebit)
	totalCredit = RoundMoney(totalCredit)
	if math.Abs(totalDebit-totalCredit) > journalBalanceTolerance {
		return 0, 0, fmt.Errorf("journal is not balanced: debit %.2f != credit %.2f", totalDebit, totalCredit)
	}
	return totalDebit, totalCredit, nil
}

package enum

import "database/sql/driver"

// AccountType is the standard Chart of Accounts classification.
type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeEquity    AccountType = "EQUITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeExpense   AccountType = "EXPENSE"
)

func (t AccountType) IsValid() bool {
	switch t {
	case AccountTypeAsset, AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue, AccountTypeExpense:
		return true
	default:
		return false
	}
}

func (t *AccountType) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*t = AccountType(v)
	return nil
}

func (t AccountType) Value() (driver.Value, error) {
	return string(t), nil
}

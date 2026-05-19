package enum

import "database/sql/driver"

type Ownership string

const (
	OwnershipCompany  Ownership = "COMPANY"
	OwnershipCustomer Ownership = "CUSTOMER"
	OwnershipVendor   Ownership = "VENDOR"
)

func (o Ownership) IsValid() bool {
	switch o {
	case OwnershipCompany, OwnershipCustomer, OwnershipVendor:
		return true
	default:
		return false
	}
}

func (o *Ownership) Scan(value interface{}) error {
	var v string
	if err := scanString(&v, value); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	*o = Ownership(v)
	return nil
}

func (o Ownership) Value() (driver.Value, error) {
	return string(o), nil
}

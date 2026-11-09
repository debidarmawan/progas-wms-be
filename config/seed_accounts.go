package config

import (
	"errors"
	"log"
	"progas-wms-be/enum"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type accountSeed struct {
	Code       string
	Name       string
	Type       enum.AccountType
	ParentCode string
}

// SeedAccounts installs a starter Chart of Accounts. Idempotent by code, and it
// never overwrites accounts an admin has already edited.
func SeedAccounts(db *gorm.DB) {
	accounts := []accountSeed{
		{Code: "1-1000", Name: "Aset", Type: enum.AccountTypeAsset},
		{Code: "1-1100", Name: "Kas", Type: enum.AccountTypeAsset, ParentCode: "1-1000"},
		{Code: "1-1200", Name: "Bank", Type: enum.AccountTypeAsset, ParentCode: "1-1000"},
		{Code: "1-1300", Name: "Piutang Usaha", Type: enum.AccountTypeAsset, ParentCode: "1-1000"},
		{Code: "1-1400", Name: "Persediaan Tabung", Type: enum.AccountTypeAsset, ParentCode: "1-1000"},
		{Code: "1-1500", Name: "Persediaan Gas", Type: enum.AccountTypeAsset, ParentCode: "1-1000"},
		{Code: "1-1600", Name: "PPN Masukan", Type: enum.AccountTypeAsset, ParentCode: "1-1000"},

		{Code: "2-1000", Name: "Liabilitas", Type: enum.AccountTypeLiability},
		{Code: "2-1100", Name: "Utang Usaha", Type: enum.AccountTypeLiability, ParentCode: "2-1000"},
		{Code: "2-1200", Name: "PPN Keluaran", Type: enum.AccountTypeLiability, ParentCode: "2-1000"},
		{Code: "2-1300", Name: "Uang Muka Pelanggan", Type: enum.AccountTypeLiability, ParentCode: "2-1000"},

		{Code: "3-1000", Name: "Ekuitas", Type: enum.AccountTypeEquity},
		{Code: "3-1100", Name: "Modal Disetor", Type: enum.AccountTypeEquity, ParentCode: "3-1000"},
		{Code: "3-1200", Name: "Laba Ditahan", Type: enum.AccountTypeEquity, ParentCode: "3-1000"},

		{Code: "4-1000", Name: "Pendapatan", Type: enum.AccountTypeRevenue},
		{Code: "4-1100", Name: "Pendapatan Penjualan Gas", Type: enum.AccountTypeRevenue, ParentCode: "4-1000"},
		{Code: "4-1200", Name: "Pendapatan Sewa Tabung", Type: enum.AccountTypeRevenue, ParentCode: "4-1000"},
		{Code: "4-1300", Name: "Pendapatan Pengiriman", Type: enum.AccountTypeRevenue, ParentCode: "4-1000"},

		{Code: "5-1000", Name: "Beban", Type: enum.AccountTypeExpense},
		{Code: "5-1100", Name: "Harga Pokok Penjualan", Type: enum.AccountTypeExpense, ParentCode: "5-1000"},
		{Code: "5-1200", Name: "Beban BBM & Transportasi", Type: enum.AccountTypeExpense, ParentCode: "5-1000"},
		{Code: "5-1300", Name: "Beban Pemeliharaan Tabung", Type: enum.AccountTypeExpense, ParentCode: "5-1000"},
		{Code: "5-1400", Name: "Beban Gaji", Type: enum.AccountTypeExpense, ParentCode: "5-1000"},
	}

	byCode := make(map[string]*model.Account, len(accounts))
	for _, seed := range accounts {
		account, err := ensureAccount(db, seed, byCode)
		if err != nil {
			log.Printf("Account seed: %s: %v", seed.Code, err)
			continue
		}
		byCode[seed.Code] = account
	}

	log.Println("Account seed completed")
}

func ensureAccount(db *gorm.DB, seed accountSeed, byCode map[string]*model.Account) (*model.Account, error) {
	var account model.Account
	err := db.Where("code = ?", seed.Code).First(&account).Error
	if err == nil {
		return &account, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	account = model.Account{
		Code:     seed.Code,
		Name:     seed.Name,
		Type:     seed.Type,
		IsActive: true,
	}
	if seed.ParentCode != "" {
		parent, ok := byCode[seed.ParentCode]
		if !ok {
			return nil, errors.New("parent account " + seed.ParentCode + " not seeded")
		}
		account.ParentId = parent.Id
	}
	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

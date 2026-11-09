package model

import "time"

// JournalEntry is a posted double-entry voucher. Entries are immutable: to
// correct one, post a reversing entry.
type JournalEntry struct {
	BaseModel
	EntryNumber string             `gorm:"not null;type:varchar(50);uniqueIndex"`
	EntryDate   time.Time          `gorm:"not null;index"`
	Description string             `gorm:"type:varchar(255)"`
	SourceType  string             `gorm:"not null;type:varchar(30);default:MANUAL"`
	SourceId    string             `gorm:"type:varchar(36);index"`
	TotalDebit  float64            `gorm:"type:decimal(15,2);not null;default:0"`
	TotalCredit float64            `gorm:"type:decimal(15,2);not null;default:0"`
	Lines       []JournalEntryLine `gorm:"foreignKey:JournalEntryId"`
}

type JournalEntryLine struct {
	BaseModel
	JournalEntryId string  `gorm:"not null;type:varchar(36);index"`
	AccountId      string  `gorm:"not null;type:varchar(36);index"`
	Account        Account `gorm:"foreignKey:AccountId"`
	Debit          float64 `gorm:"type:decimal(15,2);not null;default:0"`
	Credit         float64 `gorm:"type:decimal(15,2);not null;default:0"`
	Memo           string  `gorm:"type:varchar(255)"`
}

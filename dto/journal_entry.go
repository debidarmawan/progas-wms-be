package dto

import "time"

type CreateJournalLineRequest struct {
	AccountId string  `json:"account_id" validate:"required,max=36"`
	Debit     float64 `json:"debit" validate:"gte=0"`
	Credit    float64 `json:"credit" validate:"gte=0"`
	Memo      string  `json:"memo" validate:"max=255"`
}

type CreateJournalEntryRequest struct {
	EntryDate   string                     `json:"entry_date" validate:"required,datetime=2006-01-02"`
	Description string                     `json:"description" validate:"max=255"`
	Lines       []CreateJournalLineRequest `json:"lines" validate:"required,min=2,dive"`
}

type JournalLineResponse struct {
	AccountId   string  `json:"account_id"`
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Memo        string  `json:"memo"`
}

type JournalEntryResponse struct {
	Id          string                `json:"id"`
	EntryNumber string                `json:"entry_number"`
	EntryDate   time.Time             `json:"entry_date"`
	Description string                `json:"description"`
	SourceType  string                `json:"source_type"`
	SourceId    string                `json:"source_id,omitempty"`
	TotalDebit  float64               `json:"total_debit"`
	TotalCredit float64               `json:"total_credit"`
	Lines       []JournalLineResponse `json:"lines"`
}

type PaginatedJournalEntryList struct {
	Items []JournalEntryResponse `json:"items"`
	Meta  PaginationMeta         `json:"meta"`
}

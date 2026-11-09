package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToJournalEntryResponse(entry *model.JournalEntry) *dto.JournalEntryResponse {
	res := &dto.JournalEntryResponse{
		Id:          entry.Id,
		EntryNumber: entry.EntryNumber,
		EntryDate:   entry.EntryDate,
		Description: entry.Description,
		SourceType:  entry.SourceType,
		SourceId:    entry.SourceId,
		TotalDebit:  entry.TotalDebit,
		TotalCredit: entry.TotalCredit,
		Lines:       make([]dto.JournalLineResponse, 0, len(entry.Lines)),
	}
	for _, line := range entry.Lines {
		res.Lines = append(res.Lines, dto.JournalLineResponse{
			AccountId:   line.AccountId,
			AccountCode: line.Account.Code,
			AccountName: line.Account.Name,
			Debit:       line.Debit,
			Credit:      line.Credit,
			Memo:        line.Memo,
		})
	}
	return res
}

func ToJournalEntryResponses(entries []model.JournalEntry) []dto.JournalEntryResponse {
	responses := make([]dto.JournalEntryResponse, 0, len(entries))
	for i := range entries {
		responses = append(responses, *ToJournalEntryResponse(&entries[i]))
	}
	return responses
}

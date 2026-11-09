package usecase

import (
	"fmt"
	"strings"
	"time"

	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

const journalDateLayout = "2006-01-02"

type JournalEntryUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.JournalEntryResponse], global.ErrorResponse)
	FindById(id string) (*dto.JournalEntryResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateJournalEntryRequest) global.ErrorResponse
}

type journalEntryUsecase struct {
	txManager        helper.TxManager
	journalEntryRepo repository.JournalEntryRepository
	accountRepo      repository.AccountRepository
	auditLogRepo     repository.AuditLogRepository
}

func NewJournalEntryUsecase(
	txManager helper.TxManager,
	journalEntryRepo repository.JournalEntryRepository,
	accountRepo repository.AccountRepository,
	auditLogRepo repository.AuditLogRepository,
) JournalEntryUsecase {
	return &journalEntryUsecase{
		txManager:        txManager,
		journalEntryRepo: journalEntryRepo,
		accountRepo:      accountRepo,
		auditLogRepo:     auditLogRepo,
	}
}

func (u *journalEntryUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.JournalEntryResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)

	entries, total, err := u.journalEntryRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.JournalEntryResponse]{
		Items: mapper.ToJournalEntryResponses(entries),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *journalEntryUsecase) FindById(id string) (*dto.JournalEntryResponse, global.ErrorResponse) {
	entry, err := u.journalEntryRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToJournalEntryResponse(entry), nil
}

func (u *journalEntryUsecase) Create(actorUserId string, req *dto.CreateJournalEntryRequest) global.ErrorResponse {
	entryDate, parseErr := time.Parse(journalDateLayout, strings.TrimSpace(req.EntryDate))
	if parseErr != nil {
		return global.BadRequestError("entry_date must use YYYY-MM-DD format")
	}

	inputs := make([]helper.JournalLineInput, 0, len(req.Lines))
	for _, line := range req.Lines {
		inputs = append(inputs, helper.JournalLineInput{
			AccountId: strings.TrimSpace(line.AccountId),
			Debit:     line.Debit,
			Credit:    line.Credit,
		})
	}
	totalDebit, totalCredit, validationErr := helper.ValidateJournalLines(inputs)
	if validationErr != nil {
		return global.BadRequestError(validationErr.Error())
	}

	lines := make([]model.JournalEntryLine, 0, len(req.Lines))
	for i, line := range req.Lines {
		if err := u.ensureAccountUsable(line.AccountId, i); err != nil {
			return err
		}
		lines = append(lines, model.JournalEntryLine{
			AccountId: strings.TrimSpace(line.AccountId),
			Debit:     helper.RoundMoney(line.Debit),
			Credit:    helper.RoundMoney(line.Credit),
			Memo:      strings.TrimSpace(line.Memo),
		})
	}

	entry := &model.JournalEntry{
		EntryNumber: helper.GenerateJournalEntryNumber(),
		EntryDate:   entryDate,
		Description: strings.TrimSpace(req.Description),
		SourceType:  constant.JournalSourceManual,
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
		Lines:       lines,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.journalEntryRepo.Create(tx, entry); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditJournalCreate, constant.AuditObjectJournalEntry, entry.Id, map[string]any{
		"entry_number": entry.EntryNumber,
		"total_debit":  entry.TotalDebit,
		"total_credit": entry.TotalCredit,
	})

	return nil
}

func (u *journalEntryUsecase) ensureAccountUsable(accountId string, index int) global.ErrorResponse {
	account, err := u.accountRepo.FindById(strings.TrimSpace(accountId))
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return global.BadRequestError(fmt.Sprintf("line %d: account not found", index+1))
		}
		return err
	}
	if !account.IsActive {
		return global.BadRequestError(fmt.Sprintf("line %d: account %s is inactive", index+1, account.Code))
	}
	return nil
}

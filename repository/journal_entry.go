package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type JournalEntryRepository interface {
	FindAll(page, limit int, search string) ([]model.JournalEntry, int64, global.ErrorResponse)
	FindById(id string) (*model.JournalEntry, global.ErrorResponse)
	CountLinesByAccountId(accountId string) (int64, global.ErrorResponse)
	Create(tx helper.Tx, entry *model.JournalEntry) global.ErrorResponse
}

type journalEntryRepository struct {
	db *gorm.DB
}

func NewJournalEntryRepository(db *gorm.DB) JournalEntryRepository {
	return &journalEntryRepository{db: db}
}

func (r *journalEntryRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *journalEntryRepository) FindAll(page, limit int, search string) ([]model.JournalEntry, int64, global.ErrorResponse) {
	var entries []model.JournalEntry
	var total int64

	query := r.db.Model(&model.JournalEntry{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Where("entry_number LIKE ? OR description LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	if err := query.
		Preload("Lines.Account").
		Order("entry_date desc, entry_number desc").
		Offset(offset).Limit(limit).
		Find(&entries).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return entries, total, nil
}

func (r *journalEntryRepository) FindById(id string) (*model.JournalEntry, global.ErrorResponse) {
	var entry model.JournalEntry
	if err := r.db.Preload("Lines.Account").Where("id = ?", id).First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("journal entry not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &entry, nil
}

func (r *journalEntryRepository) CountLinesByAccountId(accountId string) (int64, global.ErrorResponse) {
	var total int64
	if err := r.db.Model(&model.JournalEntryLine{}).Where("account_id = ?", accountId).Count(&total).Error; err != nil {
		return 0, global.InternalServerError(err)
	}
	return total, nil
}

func (r *journalEntryRepository) Create(tx helper.Tx, entry *model.JournalEntry) global.ErrorResponse {
	db := r.dbFromTx(tx)
	lines := entry.Lines
	entry.Lines = nil
	if err := db.Create(entry).Error; err != nil {
		return global.InternalServerError(err)
	}
	for i := range lines {
		lines[i].JournalEntryId = entry.Id
	}
	if len(lines) > 0 {
		if err := db.Create(&lines).Error; err != nil {
			return global.InternalServerError(err)
		}
	}
	entry.Lines = lines
	return nil
}

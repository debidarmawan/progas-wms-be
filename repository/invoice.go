package repository

import (
	"errors"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	FindAll(page, limit int, search string) ([]model.Invoice, int64, global.ErrorResponse)
	FindById(id string) (*model.Invoice, global.ErrorResponse)
	Create(tx helper.Tx, invoice *model.Invoice) global.ErrorResponse
	AdjustPaidAmount(tx helper.Tx, invoiceId string, delta float64, status string) global.ErrorResponse
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) dbFromTx(tx helper.Tx) *gorm.DB {
	if tx != nil {
		return tx.Get()
	}
	return r.db
}

func (r *invoiceRepository) FindAll(page, limit int, search string) ([]model.Invoice, int64, global.ErrorResponse) {
	var invoices []model.Invoice
	var total int64

	query := r.db.Model(&model.Invoice{})
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		query = query.Joins("Customer").Where(
			"invoice.invoice_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?",
			pattern, pattern, pattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}

	offset := (page - 1) * limit
	findQuery := r.db.Preload("Customer").Preload("DeliveryOrder").Order("invoice.created_at desc").Offset(offset).Limit(limit)
	if helper.HasSearch(search) {
		pattern := helper.SearchPattern(search)
		findQuery = findQuery.Joins("Customer").Where(
			"invoice.invoice_number LIKE ? OR customer.name LIKE ? OR customer.code LIKE ?",
			pattern, pattern, pattern,
		)
	}

	if err := findQuery.Find(&invoices).Error; err != nil {
		return nil, 0, global.InternalServerError(err)
	}
	return invoices, total, nil
}

func (r *invoiceRepository) FindById(id string) (*model.Invoice, global.ErrorResponse) {
	var invoice model.Invoice
	err := r.db.Preload("Customer").Preload("DeliveryOrder").Preload("Payments").
		Where("id = ?", id).First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NotFoundError("Invoice not found")
		}
		return nil, global.InternalServerError(err)
	}
	return &invoice, nil
}

func (r *invoiceRepository) Create(tx helper.Tx, invoice *model.Invoice) global.ErrorResponse {
	if err := r.dbFromTx(tx).Create(invoice).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *invoiceRepository) AdjustPaidAmount(tx helper.Tx, invoiceId string, delta float64, status string) global.ErrorResponse {
	result := r.dbFromTx(tx).Model(&model.Invoice{}).
		Where("id = ?", invoiceId).
		Updates(map[string]any{
			"paid_amount": gorm.Expr("paid_amount + ?", delta),
			"status":      status,
		})
	if result.Error != nil {
		return global.InternalServerError(result.Error)
	}
	if result.RowsAffected == 0 {
		return global.NotFoundError("Invoice not found")
	}
	return nil
}

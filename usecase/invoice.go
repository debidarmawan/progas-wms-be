package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"time"
)

type InvoiceUsecase interface {
	CreateForDeliveryOrder(tx helper.Tx, order *model.DeliveryOrder, customer *model.Customer, details []model.DeliveryOrderDetail) (*model.Invoice, global.ErrorResponse)
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.InvoiceResponse], global.ErrorResponse)
	FindById(id string) (*dto.InvoiceResponse, global.ErrorResponse)
}

type invoiceUsecase struct {
	invoiceRepo  repository.InvoiceRepository
	customerRepo repository.CustomerRepository
}

func NewInvoiceUsecase(invoiceRepo repository.InvoiceRepository, customerRepo repository.CustomerRepository) InvoiceUsecase {
	return &invoiceUsecase{invoiceRepo: invoiceRepo, customerRepo: customerRepo}
}

func (u *invoiceUsecase) CreateForDeliveryOrder(tx helper.Tx, order *model.DeliveryOrder, customer *model.Customer, details []model.DeliveryOrderDetail) (*model.Invoice, global.ErrorResponse) {
	var total float64
	for _, d := range details {
		total += d.LineTotal
	}

	now := time.Now()
	invoice := &model.Invoice{
		BaseModel:       model.BaseModel{CreatedBy: order.CreatedBy},
		InvoiceNumber:   helper.GenerateInvoiceNumber(),
		DeliveryOrderId: order.Id,
		CustomerId:      customer.Id,
		InvoiceDate:     now,
		DueDate:         now.AddDate(0, 0, customer.PaymentTermDays),
		Status:          enum.InvoiceStatusUnpaid,
		TotalAmount:     total,
	}

	if err := u.invoiceRepo.Create(tx, invoice); err != nil {
		return nil, err
	}
	if err := u.customerRepo.AdjustOutstandingBalance(tx, customer.Id, total); err != nil {
		return nil, err
	}
	return invoice, nil
}

func (u *invoiceUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.InvoiceResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	invoices, total, err := u.invoiceRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.InvoiceResponse]{
		Items: mapper.ToInvoiceResponses(invoices),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *invoiceUsecase) FindById(id string) (*dto.InvoiceResponse, global.ErrorResponse) {
	invoice, err := u.invoiceRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToInvoiceResponse(invoice, true), nil
}

package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"time"
)

type CustomerPOUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerPOResponse], global.ErrorResponse)
	FindById(id string) (*dto.CustomerPOResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateCustomerPORequest) (*dto.CustomerPOResponse, global.ErrorResponse)
	Confirm(actorUserId, id string) (*dto.CustomerPOResponse, global.ErrorResponse)
}

type customerPOUsecase struct {
	txManager      helper.TxManager
	poRepo         repository.CustomerPORepository
	customerRepo   repository.CustomerRepository
	masterItemRepo repository.MasterItemRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewCustomerPOUsecase(
	txManager helper.TxManager,
	poRepo repository.CustomerPORepository,
	customerRepo repository.CustomerRepository,
	masterItemRepo repository.MasterItemRepository,
	auditLogRepo repository.AuditLogRepository,
) CustomerPOUsecase {
	return &customerPOUsecase{
		txManager: txManager, poRepo: poRepo, customerRepo: customerRepo,
		masterItemRepo: masterItemRepo, auditLogRepo: auditLogRepo,
	}
}

func (u *customerPOUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerPOResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	pos, total, err := u.poRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.CustomerPOResponse]{
		Items: mapper.ToCustomerPOResponses(pos),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *customerPOUsecase) FindById(id string) (*dto.CustomerPOResponse, global.ErrorResponse) {
	po, err := u.poRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToCustomerPOResponse(po, true), nil
}

func parsePODate(value string, required bool) (time.Time, *time.Time, global.ErrorResponse) {
	poDate, err := time.Parse("2006-01-02", value)
	if err != nil && required {
		return time.Time{}, nil, global.BadRequestError("po_date must use YYYY-MM-DD format")
	}
	if !required && value == "" {
		return time.Time{}, nil, nil
	}
	if err != nil {
		return time.Time{}, nil, global.BadRequestError("valid_until must use YYYY-MM-DD format")
	}
	if required {
		return poDate, nil, nil
	}
	return time.Time{}, &poDate, nil
}

func (u *customerPOUsecase) Create(actorUserId string, req *dto.CreateCustomerPORequest) (*dto.CustomerPOResponse, global.ErrorResponse) {
	customer, err := u.customerRepo.FindById(req.CustomerId)
	if err != nil {
		return nil, err
	}
	if !customer.IsActive {
		return nil, global.BadRequestError("customer is not active")
	}

	poDate, _, parseErr := parsePODate(req.PODate, true)
	if parseErr != nil {
		return nil, parseErr
	}
	_, validUntil, parseErr := parsePODate(req.ValidUntil, false)
	if parseErr != nil {
		return nil, parseErr
	}

	seenItems := make(map[string]struct{}, len(req.Lines))
	lines := make([]model.CustomerPOLine, 0, len(req.Lines))
	for _, lineReq := range req.Lines {
		if _, exists := seenItems[lineReq.MasterItemId]; exists {
			return nil, global.BadRequestError("customer PO cannot contain duplicate master items")
		}
		item, itemErr := u.masterItemRepo.FindById(lineReq.MasterItemId)
		if itemErr != nil {
			return nil, itemErr
		}
		seenItems[lineReq.MasterItemId] = struct{}{}
		lines = append(lines, model.CustomerPOLine{
			MasterItemId: item.Id,
			Quantity:     lineReq.Quantity,
		})
	}

	po := &model.CustomerPO{
		BaseModel:   model.BaseModel{CreatedBy: actorUserId},
		PONumber:    req.PONumber,
		CustomerId:  customer.Id,
		PODate:      poDate,
		ValidUntil:  validUntil,
		DocumentURL: req.DocumentURL,
		Status:      enum.CustomerPOStatusDraft,
		Lines:       lines,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()
	if err := u.poRepo.Create(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerPOCreate, constant.AuditObjectCustomerPO, po.Id, map[string]any{
		"po_number": po.PONumber, "customer_id": po.CustomerId,
	})
	po.Customer = *customer
	return mapper.ToCustomerPOResponse(po, true), nil
}

func (u *customerPOUsecase) Confirm(actorUserId, id string) (*dto.CustomerPOResponse, global.ErrorResponse) {
	po, err := u.poRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	if po.Status != enum.CustomerPOStatusDraft {
		return nil, global.BadRequestError("only draft customer PO can be confirmed")
	}

	po.Status = enum.CustomerPOStatusConfirmed
	tx := u.txManager.New()
	defer tx.CheckPanic()
	if err := u.poRepo.Update(tx, po); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerPOConfirm, constant.AuditObjectCustomerPO, po.Id, map[string]any{
		"po_number": po.PONumber,
	})
	return mapper.ToCustomerPOResponse(po, true), nil
}

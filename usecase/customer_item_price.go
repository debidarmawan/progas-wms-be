package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"time"

	"github.com/gofiber/fiber/v3"
)

const customerItemPriceTimeFormat = time.RFC3339
const customerItemPriceDateFormat = "2006-01-02"

type CustomerItemPriceUsecase interface {
	FindAllByCustomer(customerId string, query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerItemPriceResponse], global.ErrorResponse)
	FindActive(customerId, masterItemId string) (*dto.CustomerItemPriceResponse, global.ErrorResponse)
	Create(actorUserId, customerId string, req *dto.CreateCustomerItemPriceRequest) global.ErrorResponse
	Update(actorUserId, customerId, id string, req *dto.UpdateCustomerItemPriceRequest) global.ErrorResponse
	Delete(actorUserId, customerId, id string) global.ErrorResponse
	ResolvePrice(customerId, masterItemId string, at time.Time) (float64, string, global.ErrorResponse)
}

type customerItemPriceUsecase struct {
	txManager      helper.TxManager
	priceRepo      repository.CustomerItemPriceRepository
	customerRepo   repository.CustomerRepository
	masterItemRepo repository.MasterItemRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewCustomerItemPriceUsecase(
	txManager helper.TxManager,
	priceRepo repository.CustomerItemPriceRepository,
	customerRepo repository.CustomerRepository,
	masterItemRepo repository.MasterItemRepository,
	auditLogRepo repository.AuditLogRepository,
) CustomerItemPriceUsecase {
	return &customerItemPriceUsecase{
		txManager:      txManager,
		priceRepo:      priceRepo,
		customerRepo:   customerRepo,
		masterItemRepo: masterItemRepo,
		auditLogRepo:   auditLogRepo,
	}
}

func (u *customerItemPriceUsecase) FindAllByCustomer(customerId string, query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerItemPriceResponse], global.ErrorResponse) {
	if _, err := u.customerRepo.FindById(customerId); err != nil {
		return nil, err
	}
	page, limit, _ := helper.NormalizePagination(query)
	prices, total, err := u.priceRepo.FindAllByCustomer(customerId, page, limit, helper.NormalizeSearch(query.Search))
	if err != nil {
		return nil, err
	}
	responses := make([]dto.CustomerItemPriceResponse, 0, len(prices))
	for i := range prices {
		item, itemErr := u.masterItemRepo.FindById(prices[i].MasterItemId)
		if itemErr != nil {
			return nil, itemErr
		}
		responses = append(responses, *mapper.ToCustomerItemPriceResponse(&prices[i], item, time.Now()))
	}
	return &dto.PaginatedResponse[dto.CustomerItemPriceResponse]{
		Items: responses,
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *customerItemPriceUsecase) FindActive(customerId, masterItemId string) (*dto.CustomerItemPriceResponse, global.ErrorResponse) {
	if _, err := u.customerRepo.FindById(customerId); err != nil {
		return nil, err
	}
	item, err := u.masterItemRepo.FindById(masterItemId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	price, priceErr := u.priceRepo.FindActive(customerId, masterItemId, now)
	if priceErr == nil {
		return mapper.ToCustomerItemPriceResponse(price, item, now), nil
	}
	if priceErr.GetCode() != fiber.StatusNotFound {
		return nil, priceErr
	}
	return &dto.CustomerItemPriceResponse{
		CustomerId:    customerId,
		MasterItemId:  masterItemId,
		ItemName:      item.Name,
		ItemSKU:       item.SKU,
		Price:         item.HnaPrice,
		EffectiveFrom: now.Format(customerItemPriceTimeFormat),
		IsActive:      true,
		IsFallback:    true,
	}, nil
}

func (u *customerItemPriceUsecase) ResolvePrice(customerId, masterItemId string, at time.Time) (float64, string, global.ErrorResponse) {
	if _, err := u.customerRepo.FindById(customerId); err != nil {
		return 0, "", err
	}
	item, err := u.masterItemRepo.FindById(masterItemId)
	if err != nil {
		return 0, "", err
	}
	price, priceErr := u.priceRepo.FindActive(customerId, masterItemId, at)
	if priceErr == nil {
		return price.Price, "CUSTOM", nil
	}
	if priceErr.GetCode() != fiber.StatusNotFound {
		return 0, "", priceErr
	}
	return item.HnaPrice, "HNA", nil
}

func (u *customerItemPriceUsecase) Create(actorUserId, customerId string, req *dto.CreateCustomerItemPriceRequest) global.ErrorResponse {
	if _, err := u.customerRepo.FindById(customerId); err != nil {
		return err
	}
	if _, err := u.masterItemRepo.FindById(req.MasterItemId); err != nil {
		return err
	}
	from, to, err := parsePricePeriod(req.EffectiveFrom, req.EffectiveTo)
	if err != nil {
		return global.BadRequestError(err.Error())
	}
	overlap, errResponse := u.priceRepo.HasOverlap(customerId, req.MasterItemId, from, to, "")
	if errResponse != nil {
		return errResponse
	}
	if overlap {
		return global.BadRequestError("price period overlaps an existing customer item price")
	}
	price := &model.CustomerItemPrice{CustomerId: customerId, MasterItemId: req.MasterItemId, Price: req.Price, EffectiveFrom: from, EffectiveTo: to}
	price.CreatedBy = actorUserId
	tx := u.txManager.New()
	defer tx.CheckPanic()
	if errResponse = u.priceRepo.Create(tx, price); errResponse != nil {
		tx.Rollback()
		return errResponse
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}
	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerPriceCreate, constant.AuditObjectCustomerPrice, price.Id, map[string]any{"customer_id": customerId, "master_item_id": req.MasterItemId, "price": req.Price})
	return nil
}

func (u *customerItemPriceUsecase) Update(actorUserId, customerId, id string, req *dto.UpdateCustomerItemPriceRequest) global.ErrorResponse {
	price, err := u.priceRepo.FindById(id)
	if err != nil {
		return err
	}
	if price.CustomerId != customerId {
		return global.NotFoundError("Customer item price not found")
	}
	if !price.EffectiveFrom.After(time.Now()) {
		return global.BadRequestError("prices that already started cannot be changed")
	}
	from, to, parseErr := parsePricePeriod(req.EffectiveFrom, req.EffectiveTo)
	if parseErr != nil {
		return global.BadRequestError(parseErr.Error())
	}
	overlap, errResponse := u.priceRepo.HasOverlap(customerId, price.MasterItemId, from, to, price.Id)
	if errResponse != nil {
		return errResponse
	}
	if overlap {
		return global.BadRequestError("price period overlaps an existing customer item price")
	}
	price.Price = req.Price
	price.EffectiveFrom = from
	price.EffectiveTo = to
	price.UpdatedBy = actorUserId
	tx := u.txManager.New()
	defer tx.CheckPanic()
	if errResponse = u.priceRepo.Update(tx, price); errResponse != nil {
		tx.Rollback()
		return errResponse
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}
	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerPriceUpdate, constant.AuditObjectCustomerPrice, price.Id, map[string]any{"customer_id": customerId, "master_item_id": price.MasterItemId, "price": req.Price})
	return nil
}

func (u *customerItemPriceUsecase) Delete(actorUserId, customerId, id string) global.ErrorResponse {
	price, err := u.priceRepo.FindById(id)
	if err != nil {
		return err
	}
	if price.CustomerId != customerId {
		return global.NotFoundError("Customer item price not found")
	}
	price.DeletedBy = actorUserId
	tx := u.txManager.New()
	defer tx.CheckPanic()
	if errResponse := u.priceRepo.Delete(tx, price); errResponse != nil {
		tx.Rollback()
		return errResponse
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}
	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerPriceDelete, constant.AuditObjectCustomerPrice, price.Id, map[string]any{"customer_id": customerId, "master_item_id": price.MasterItemId})
	return nil
}

func parsePricePeriod(fromValue string, toValue *string) (time.Time, *time.Time, error) {
	fromDate, err := time.ParseInLocation(customerItemPriceDateFormat, fromValue, time.Local)
	if err != nil {
		return time.Time{}, nil, fiber.NewError(fiber.StatusBadRequest, "effective_from must be a date in YYYY-MM-DD format")
	}
	from := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)
	if toValue == nil || *toValue == "" {
		return from, nil, nil
	}
	toDate, err := time.ParseInLocation(customerItemPriceDateFormat, *toValue, time.Local)
	if err != nil {
		return time.Time{}, nil, fiber.NewError(fiber.StatusBadRequest, "effective_to must be a date in YYYY-MM-DD format")
	}
	to := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), time.Local)
	if !to.After(from) {
		return time.Time{}, nil, fiber.NewError(fiber.StatusBadRequest, "effective_to must be after effective_from")
	}
	return from, &to, nil
}

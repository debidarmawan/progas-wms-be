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
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type CylinderUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CylinderResponse], global.ErrorResponse)
	FindById(id string) (*dto.CylinderResponse, global.ErrorResponse)
	FindByBarcode(barcode string) (*dto.CylinderResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateCylinderRequest) global.ErrorResponse
	Update(actorUserId string, id string, req *dto.UpdateCylinderRequest) global.ErrorResponse
	History(id string) (*dto.CylinderHistoryResponse, global.ErrorResponse)
}

type cylinderUsecase struct {
	txManager      helper.TxManager
	cylinderRepo   repository.CylinderRepository
	masterItemRepo repository.MasterItemRepository
	customerRepo   repository.CustomerRepository
	vendorRepo     repository.VendorRepository
	ledgerRepo     repository.CylinderLedgerRepository
	userRepo       repository.UserRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewCylinderUsecase(
	txManager helper.TxManager,
	cylinderRepo repository.CylinderRepository,
	masterItemRepo repository.MasterItemRepository,
	customerRepo repository.CustomerRepository,
	vendorRepo repository.VendorRepository,
	ledgerRepo repository.CylinderLedgerRepository,
	userRepo repository.UserRepository,
	auditLogRepo repository.AuditLogRepository,
) CylinderUsecase {
	return &cylinderUsecase{
		txManager:      txManager,
		cylinderRepo:   cylinderRepo,
		masterItemRepo: masterItemRepo,
		customerRepo:   customerRepo,
		vendorRepo:     vendorRepo,
		ledgerRepo:     ledgerRepo,
		userRepo:       userRepo,
		auditLogRepo:   auditLogRepo,
	}
}

func (u *cylinderUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CylinderResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	cylinders, total, err := u.cylinderRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.CylinderResponse]{
		Items: mapper.ToCylinderResponses(cylinders),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *cylinderUsecase) FindById(id string) (*dto.CylinderResponse, global.ErrorResponse) {
	cylinder, err := u.cylinderRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToCylinderResponse(cylinder), nil
}

func (u *cylinderUsecase) FindByBarcode(barcode string) (*dto.CylinderResponse, global.ErrorResponse) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return nil, global.BadRequestError("barcode is required")
	}
	cylinder, err := u.cylinderRepo.FindByBarcode(barcode)
	if err != nil {
		return nil, err
	}
	return mapper.ToCylinderResponse(cylinder), nil
}

func (u *cylinderUsecase) Create(actorUserId string, req *dto.CreateCylinderRequest) global.ErrorResponse {
	existing, err := u.cylinderRepo.FindByBarcode(req.BarcodeSN)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existing != nil {
		return global.BadRequestError("barcode already registered")
	}

	item, err := u.masterItemRepo.FindById(req.ItemId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return global.BadRequestError("invalid item")
		}
		return err
	}
	if !item.IsSerialized {
		return global.BadRequestError("item is not serialized; cannot register cylinder")
	}

	ownership := enum.Ownership(req.OwnershipType)
	ownerName := ""
	var ownerId *string
	if req.OwnerId != "" {
		ownerId = &req.OwnerId
	}
	if !helper.ValidateOwnership(ownership, ownerId) {
		return global.BadRequestError("invalid ownership: CUSTOMER and VENDOR require owner_id; COMPANY must not have owner_id")
	}
	if ownership == enum.OwnershipCustomer {
		customer, err := u.customerRepo.FindById(*ownerId)
		if err != nil {
			if err.GetCode() == fiber.StatusNotFound {
				return global.BadRequestError("invalid customer owner_id")
			}
			return err
		}
		ownerName = customer.Name
	}
	if ownership == enum.OwnershipVendor {
		vendor, err := u.vendorRepo.FindById(*ownerId)
		if err != nil {
			if err.GetCode() == fiber.StatusNotFound {
				return global.BadRequestError("invalid vendor owner_id")
			}
			return err
		}
		if !vendor.IsActive {
			return global.BadRequestError("vendor is not active")
		}
		ownerName = vendor.Name
	}

	hydrotestDate, parseErr := time.Parse(time.RFC3339, req.LastHydrotestDate)
	if parseErr != nil {
		hydrotestDate, parseErr = time.Parse("2006-01-02", req.LastHydrotestDate)
	}
	if parseErr != nil {
		return global.BadRequestError("invalid last_hydrotest_date format (use YYYY-MM-DD or RFC3339)")
	}
	if !helper.ValidateHydrotestDate(hydrotestDate) {
		return global.BadRequestError("last hydrotest date is invalid or expired")
	}

	cylinder := &model.Cylinder{
		BarcodeSN:         req.BarcodeSN,
		ItemId:            req.ItemId,
		OwnershipType:     ownership,
		OwnerId:           ownerId,
		Status:            enum.CylinderStatusEmpty,
		LastHydrotestDate: hydrotestDate,
		Remarks:           req.Remarks,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.cylinderRepo.Create(tx, cylinder); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	changes := []dto.CylinderHistoryChange{
		{Field: "barcode_sn", Label: "Barcode", New: cylinder.BarcodeSN},
		{Field: "item_id", Label: "Produk Gas", New: item.Name},
		{Field: "ownership_type", Label: "Kepemilikan", New: string(cylinder.OwnershipType)},
		{Field: "owner_id", Label: "Pemilik", New: ownerName},
		{Field: "last_hydrotest_date", Label: "Tanggal Hydrotest", New: cylinder.LastHydrotestDate.Format("2006-01-02")},
		{Field: "remarks", Label: "Remarks", New: cylinder.Remarks},
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCylinderCreate, constant.AuditObjectCylinder, cylinder.Id, map[string]any{
		"changes": changes,
	})

	return nil
}

func (u *cylinderUsecase) Update(actorUserId string, id string, req *dto.UpdateCylinderRequest) global.ErrorResponse {
	cylinder, err := u.cylinderRepo.FindById(id)
	if err != nil {
		return err
	}

	if cylinder.Status != enum.CylinderStatusEmpty && cylinder.Status != enum.CylinderStatusReadyToFill {
		return global.BadRequestError("tabung hanya bisa diedit saat status EMPTY atau READY_TO_FILL")
	}

	item, err := u.masterItemRepo.FindById(req.ItemId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return global.BadRequestError("invalid item")
		}
		return err
	}
	if !item.IsSerialized {
		return global.BadRequestError("item is not serialized; cannot assign to cylinder")
	}

	ownership := enum.Ownership(req.OwnershipType)
	ownerName := ""
	var ownerId *string
	if req.OwnerId != "" {
		ownerId = &req.OwnerId
	}
	if !helper.ValidateOwnership(ownership, ownerId) {
		return global.BadRequestError("invalid ownership: CUSTOMER and VENDOR require owner_id; COMPANY must not have owner_id")
	}
	if ownership == enum.OwnershipCustomer {
		customer, err := u.customerRepo.FindById(*ownerId)
		if err != nil {
			if err.GetCode() == fiber.StatusNotFound {
				return global.BadRequestError("invalid customer owner_id")
			}
			return err
		}
		ownerName = customer.Name
	}
	if ownership == enum.OwnershipVendor {
		vendor, err := u.vendorRepo.FindById(*ownerId)
		if err != nil {
			if err.GetCode() == fiber.StatusNotFound {
				return global.BadRequestError("invalid vendor owner_id")
			}
			return err
		}
		if !vendor.IsActive {
			return global.BadRequestError("vendor is not active")
		}
		ownerName = vendor.Name
	}

	oldOwnerName := ""
	if cylinder.OwnerId != nil && *cylinder.OwnerId != "" {
		if cylinder.OwnershipType == enum.OwnershipCustomer {
			if customer, err := u.customerRepo.FindById(*cylinder.OwnerId); err == nil {
				oldOwnerName = customer.Name
			}
		} else if cylinder.OwnershipType == enum.OwnershipVendor {
			if vendor, err := u.vendorRepo.FindById(*cylinder.OwnerId); err == nil {
				oldOwnerName = vendor.Name
			}
		}
	}

	hydrotestDate, parseErr := time.Parse(time.RFC3339, req.LastHydrotestDate)
	if parseErr != nil {
		hydrotestDate, parseErr = time.Parse("2006-01-02", req.LastHydrotestDate)
	}
	if parseErr != nil {
		return global.BadRequestError("invalid last_hydrotest_date format (use YYYY-MM-DD or RFC3339)")
	}
	if !helper.ValidateHydrotestDate(hydrotestDate) {
		return global.BadRequestError("last hydrotest date is invalid or expired")
	}

	oldItemName := ""
	if cylinder.MasterItem.Id != "" {
		oldItemName = cylinder.MasterItem.Name
	}

	changes := []dto.CylinderHistoryChange{}
	appendCylinderChange(&changes, "barcode_sn", "Barcode", cylinder.BarcodeSN, req.BarcodeSN)
	appendCylinderChange(&changes, "item_id", "Produk Gas", oldItemName, item.Name)
	appendCylinderChange(&changes, "ownership_type", "Kepemilikan", string(cylinder.OwnershipType), string(ownership))
	appendCylinderChange(&changes, "owner_id", "Pemilik", oldOwnerName, ownerName)
	appendCylinderChange(&changes, "last_hydrotest_date", "Tanggal Hydrotest", cylinder.LastHydrotestDate.Format("2006-01-02"), hydrotestDate.Format("2006-01-02"))
	appendCylinderChange(&changes, "remarks", "Remarks", cylinder.Remarks, req.Remarks)

	cylinder.BarcodeSN = req.BarcodeSN
	cylinder.ItemId = req.ItemId
	cylinder.OwnershipType = ownership
	cylinder.OwnerId = ownerId
	cylinder.LastHydrotestDate = hydrotestDate
	cylinder.Remarks = req.Remarks

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.cylinderRepo.Update(tx, cylinder); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCylinderUpdate, constant.AuditObjectCylinder, cylinder.Id, map[string]any{
		"changes": changes,
	})

	return nil
}

func appendCylinderChange(changes *[]dto.CylinderHistoryChange, field, label, oldV, newV string) {
	if strings.TrimSpace(oldV) == strings.TrimSpace(newV) {
		return
	}
	*changes = append(*changes, dto.CylinderHistoryChange{Field: field, Label: label, Old: oldV, New: newV})
}

func (u *cylinderUsecase) History(id string) (*dto.CylinderHistoryResponse, global.ErrorResponse) {
	cylinder, err := u.cylinderRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	auditLogs, err := u.auditLogRepo.FindByObject(constant.AuditObjectCylinder, cylinder.Id)
	if err != nil {
		return nil, err
	}
	ledgerEntries, err := u.ledgerRepo.FindByCylinderId(cylinder.Id)
	if err != nil {
		return nil, err
	}

	userIds := map[string]bool{}
	entries := make([]dto.CylinderHistoryEntry, 0, len(auditLogs)+len(ledgerEntries))
	itemName := func(id string) string {
		if item, err := u.masterItemRepo.FindById(id); err == nil {
			return item.Name
		}
		return ""
	}
	for _, log := range auditLogs {
		if log.Action != constant.AuditCylinderCreate && log.Action != constant.AuditCylinderUpdate {
			continue
		}
		userIds[log.UserId] = true
		entries = append(entries, dto.CylinderHistoryEntry{
			Id:          log.Id,
			Action:      log.Action,
			ActionLabel: mapper.CylinderHistoryActionLabel(log.Action),
			UserId:      log.UserId,
			CreatedAt:   log.CreatedAt.Format(time.RFC3339),
			Changes:     mapper.ParseCylinderHistoryChanges(log.Details, itemName),
		})
	}

	userNames := map[string]string{}
	for userId := range userIds {
		if user, err := u.userRepo.FindById(userId); err == nil {
			userNames[userId] = user.Name
		}
	}
	for i := range entries {
		entries[i].UserName = userNames[entries[i].UserId]
	}

	for _, entry := range ledgerEntries {
		entries = append(entries, dto.CylinderHistoryEntry{
			Id:          entry.Id,
			Action:      entry.Action,
			ActionLabel: mapper.CylinderHistoryActionLabel(entry.Action),
			CreatedAt:   entry.CreatedAt.Format(time.RFC3339),
			Changes: []dto.CylinderHistoryChange{
				{Field: "status", Label: "Status", Old: string(entry.FromStatus), New: string(entry.ToStatus)},
			},
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].CreatedAt < entries[j].CreatedAt
	})

	return mapper.ToCylinderHistoryResponse(cylinder, entries), nil
}

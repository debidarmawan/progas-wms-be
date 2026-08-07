package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
)

type DashboardUsecase interface {
	GetSummary() (*dto.DashboardSummaryResponse, global.ErrorResponse)
}

type dashboardUsecase struct {
	dashboardRepo      repository.DashboardRepository
	sparepartStockRepo repository.SparepartStockRepository
	cylinderRepo       repository.CylinderRepository
	customerRepo       repository.CustomerRepository
}

func NewDashboardUsecase(
	dashboardRepo repository.DashboardRepository,
	sparepartStockRepo repository.SparepartStockRepository,
	cylinderRepo repository.CylinderRepository,
	customerRepo repository.CustomerRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		dashboardRepo:      dashboardRepo,
		sparepartStockRepo: sparepartStockRepo,
		cylinderRepo:       cylinderRepo,
		customerRepo:       customerRepo,
	}
}

func (u *dashboardUsecase) GetSummary() (*dto.DashboardSummaryResponse, global.ErrorResponse) {
	byStatus, err := u.dashboardRepo.CountCylindersByStatus()
	if err != nil {
		return nil, err
	}

	totalOutstanding, err := u.dashboardRepo.CountOutstandingCylinders()
	if err != nil {
		return nil, err
	}

	overQuota, err := u.dashboardRepo.FindCustomersOverQuota()
	if err != nil {
		return nil, err
	}

	lowStocks, err := u.sparepartStockRepo.FindAllLowStock()
	if err != nil {
		return nil, err
	}

	expiredCount, err := u.dashboardRepo.CountHydrotestExpired()
	if err != nil {
		return nil, err
	}

	dueSoonCount, err := u.dashboardRepo.CountHydrotestDueSoon(30)
	if err != nil {
		return nil, err
	}

	lowStockAlerts := make([]dto.LowStockSparepartAlert, 0, len(lowStocks))
	for _, stock := range lowStocks {
		lowStockAlerts = append(lowStockAlerts, dto.LowStockSparepartAlert{
			ItemId:   stock.ItemId,
			ItemName: stock.MasterItem.Name,
			SKU:      stock.MasterItem.SKU,
			Quantity: stock.Quantity,
			MinStock: stock.MasterItem.MinStockAlert,
		})
	}

	quotaAlerts := make([]dto.CustomerQuotaAlert, 0, len(overQuota))
	for _, c := range overQuota {
		quotaAlerts = append(quotaAlerts, dto.CustomerQuotaAlert{
			CustomerId:       c.Id,
			CustomerCode:     c.Code,
			CustomerName:     c.Name,
			OutstandingCount: c.OutstandingCount,
			QuotaLimit:       c.CylinderQuotaLimit,
		})
	}

	overdueAlerts := make([]dto.OverdueCylinderAlert, 0)
	if err := u.buildOverdueAlerts(&overdueAlerts); err != nil {
		return nil, err
	}

	return &dto.DashboardSummaryResponse{
		CylindersByStatus:         byStatus,
		LowStockSpareparts:        lowStockAlerts,
		TotalOutstandingCylinders: int(totalOutstanding),
		CustomersOverQuota:        quotaAlerts,
		OverdueCylindersCount:     len(overdueAlerts),
		OverdueCylinders:          overdueAlerts,
		HydrotestExpiredCount:     int(expiredCount),
		HydrotestDueSoonCount:     int(dueSoonCount),
	}, nil
}

func (u *dashboardUsecase) buildOverdueAlerts(alerts *[]dto.OverdueCylinderAlert) global.ErrorResponse {
	cylinders, err := u.cylinderRepo.FindOutstanding()
	if err != nil {
		return err
	}

	customers, _, err := u.customerRepo.FindAll(1, 1000, "")
	if err != nil {
		return err
	}
	customerById := make(map[string]model.Customer, len(customers))
	for _, c := range customers {
		customerById[c.Id] = c
	}

	for _, cyl := range cylinders {
		maxDays := cyl.MasterItem.MaxDaysAtCustomer
		if !helper.IsOverdueAtCustomer(cyl.OutstandingSince, maxDays) {
			continue
		}
		if cyl.OwnerId == nil {
			continue
		}
		customer, ok := customerById[*cyl.OwnerId]
		if !ok {
			continue
		}
		*alerts = append(*alerts, dto.OverdueCylinderAlert{
			CustomerId:     customer.Id,
			CustomerCode:   customer.Code,
			CustomerName:   customer.Name,
			BarcodeSN:      cyl.BarcodeSN,
			DaysAtCustomer: helper.DaysAtCustomer(cyl.OutstandingSince),
			MaxDays:        maxDays,
		})
	}
	return nil
}

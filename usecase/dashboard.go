package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/repository"
)

type DashboardUsecase interface {
	GetSummary() (*dto.DashboardSummaryResponse, global.ErrorResponse)
}

type dashboardUsecase struct {
	dashboardRepo      repository.DashboardRepository
	sparepartStockRepo repository.SparepartStockRepository
}

func NewDashboardUsecase(
	dashboardRepo repository.DashboardRepository,
	sparepartStockRepo repository.SparepartStockRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		dashboardRepo:      dashboardRepo,
		sparepartStockRepo: sparepartStockRepo,
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

	return &dto.DashboardSummaryResponse{
		CylindersByStatus:         byStatus,
		LowStockSpareparts:        lowStockAlerts,
		TotalOutstandingCylinders: int(totalOutstanding),
		CustomersOverQuota:        quotaAlerts,
		HydrotestExpiredCount:     int(expiredCount),
		HydrotestDueSoonCount:     int(dueSoonCount),
	}, nil
}

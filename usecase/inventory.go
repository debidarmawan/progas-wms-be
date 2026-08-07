package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/repository"
)

type InventoryUsecase interface {
	VirtualWarehouse() (*dto.VirtualWarehouseResponse, global.ErrorResponse)
}

type inventoryUsecase struct {
	customerRepo repository.CustomerRepository
	cylinderRepo repository.CylinderRepository
}

func NewInventoryUsecase(
	customerRepo repository.CustomerRepository,
	cylinderRepo repository.CylinderRepository,
) InventoryUsecase {
	return &inventoryUsecase{
		customerRepo: customerRepo,
		cylinderRepo: cylinderRepo,
	}
}

func (u *inventoryUsecase) VirtualWarehouse() (*dto.VirtualWarehouseResponse, global.ErrorResponse) {
	customers, _, err := u.customerRepo.FindAll(1, 1000, "")
	if err != nil {
		return nil, err
	}

	grouped, err := u.cylinderRepo.FindOutstandingGroupedByCustomer()
	if err != nil {
		return nil, err
	}

	result := make([]dto.VirtualWarehouseCustomer, 0)
	for _, customer := range customers {
		if customer.OutstandingCount <= 0 {
			continue
		}
		cylinders := make([]dto.VirtualWarehouseCylinder, 0)
		overdueCount := 0
		if cyls, ok := grouped[customer.Id]; ok {
			for _, cyl := range cyls {
				maxDays := cyl.MasterItem.MaxDaysAtCustomer
				days := helper.DaysAtCustomer(cyl.OutstandingSince)
				isOverdue := helper.IsOverdueAtCustomer(cyl.OutstandingSince, maxDays)
				if isOverdue {
					overdueCount++
				}
				cylinders = append(cylinders, dto.VirtualWarehouseCylinder{
					BarcodeSN:      cyl.BarcodeSN,
					DaysAtCustomer: days,
					MaxDays:        maxDays,
					IsOverdue:      isOverdue,
				})
			}
		}
		result = append(result, dto.VirtualWarehouseCustomer{
			CustomerId:       customer.Id,
			CustomerCode:     customer.Code,
			CustomerName:     customer.Name,
			OutstandingCount: customer.OutstandingCount,
			OverdueCount:     overdueCount,
			Cylinders:        cylinders,
		})
	}

	return &dto.VirtualWarehouseResponse{Customers: result}, nil
}

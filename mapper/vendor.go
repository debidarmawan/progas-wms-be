package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToVendorResponse(v *model.Vendor, cylinderCount int) *dto.VendorResponse {
	res := &dto.VendorResponse{
		Id:            v.Id,
		Code:          v.Code,
		Name:          v.Name,
		ContactPerson: v.ContactPerson,
		Phone:         v.Phone,
		Email:         v.Email,
		Address:       v.Address,
		Notes:         v.Notes,
		IsActive:      v.IsActive,
		CylinderCount: cylinderCount,
	}
	if v.ContractStartDate != nil {
		res.ContractStartDate = v.ContractStartDate.Format("2006-01-02")
	}
	if v.ContractEndDate != nil {
		res.ContractEndDate = v.ContractEndDate.Format("2006-01-02")
	}
	return res
}

func ToVendorResponses(vendors []model.Vendor, counts map[string]int64) []dto.VendorResponse {
	responses := make([]dto.VendorResponse, 0, len(vendors))
	for i := range vendors {
		count := int(counts[vendors[i].Id])
		responses = append(responses, *ToVendorResponse(&vendors[i], count))
	}
	return responses
}

func ToVendorCylinderSummaries(cylinders []model.Cylinder) []dto.VendorCylinderSummary {
	summaries := make([]dto.VendorCylinderSummary, 0, len(cylinders))
	for _, cyl := range cylinders {
		itemName, gasType := "", ""
		if cyl.MasterItem.Id != "" {
			itemName = cyl.MasterItem.Name
			gasType = cyl.MasterItem.GasType
		}
		summaries = append(summaries, dto.VendorCylinderSummary{
			Id:        cyl.Id,
			BarcodeSN: cyl.BarcodeSN,
			Status:    string(cyl.Status),
			GasType:   gasType,
			ItemName:  itemName,
		})
	}
	return summaries
}

func CylindersByStatus(cylinders []model.Cylinder) map[string]int {
	counts := make(map[string]int)
	for _, cyl := range cylinders {
		counts[string(cyl.Status)]++
	}
	return counts
}

func ParseOptionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

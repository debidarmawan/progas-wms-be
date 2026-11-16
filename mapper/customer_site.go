package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToCustomerSiteResponse(site *model.CustomerSite) *dto.CustomerSiteResponse {
	return &dto.CustomerSiteResponse{
		Id:         site.Id,
		CustomerId: site.CustomerId,
		Name:       site.Name,
		Pic:        site.Pic,
		Phone:      site.Phone,
		Address:    site.Address,
		PostalCode: site.PostalCode,
		Notes:      site.Notes,
		IsPrimary:  site.IsPrimary,
		IsActive:   site.IsActive,
	}
}

func ToCustomerSiteResponses(sites []model.CustomerSite) []dto.CustomerSiteResponse {
	responses := make([]dto.CustomerSiteResponse, 0, len(sites))
	for i := range sites {
		responses = append(responses, *ToCustomerSiteResponse(&sites[i]))
	}
	return responses
}

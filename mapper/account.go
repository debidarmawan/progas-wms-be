package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/helper"
	"progas-wms-be/model"
)

func ToAccountResponse(account *model.Account) *dto.AccountResponse {
	res := &dto.AccountResponse{
		Id:            account.Id,
		Code:          account.Code,
		Name:          account.Name,
		Type:          string(account.Type),
		NormalBalance: helper.AccountNormalBalance(account.Type),
		ParentId:      account.ParentId,
		IsActive:      account.IsActive,
	}
	if account.Parent != nil {
		res.ParentCode = account.Parent.Code
	}
	return res
}

func ToAccountResponses(accounts []model.Account) []dto.AccountResponse {
	responses := make([]dto.AccountResponse, 0, len(accounts))
	for i := range accounts {
		responses = append(responses, *ToAccountResponse(&accounts[i]))
	}
	return responses
}

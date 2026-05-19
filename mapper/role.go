package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToRoleResponse(role *model.Role) *dto.RoleResponse {
	return &dto.RoleResponse{
		Id:   role.Id,
		Name: role.Name,
	}
}

func ToRoleResponses(roles []model.Role) []dto.RoleResponse {
	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, *ToRoleResponse(&role))
	}
	return responses
}

package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToUserListResponse(user *model.User) *dto.UserListResponse {
	roleName := ""
	if user.Role.Id != "" {
		roleName = user.Role.Name
	}
	res := &dto.UserListResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		RoleId:    user.RoleId,
		RoleName:  roleName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	if user.LastLoggedInAt != nil {
		res.LastLoggedInAt = user.LastLoggedInAt.Format(time.RFC3339)
	}
	return res
}

func ToUserListResponses(users []model.User) []dto.UserListResponse {
	responses := make([]dto.UserListResponse, 0, len(users))
	for i := range users {
		responses = append(responses, *ToUserListResponse(&users[i]))
	}
	return responses
}

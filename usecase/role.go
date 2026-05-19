package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/repository"
)

type RoleUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.RoleResponse], global.ErrorResponse)
	FindById(id string) (*dto.RoleResponse, global.ErrorResponse)
}

type roleUsecase struct {
	roleRepository repository.RoleRepository
}

func NewRoleUsecase(roleRepository repository.RoleRepository) RoleUsecase {
	return &roleUsecase{roleRepository}
}

func (u *roleUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.RoleResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	roles, total, err := u.roleRepository.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.RoleResponse]{
		Items: mapper.ToRoleResponses(roles),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *roleUsecase) FindById(id string) (*dto.RoleResponse, global.ErrorResponse) {
	role, err := u.roleRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToRoleResponse(role), nil
}

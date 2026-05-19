package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/mapper"
	"progas-wms-be/repository"
)

type RoleUsecase interface {
	FindAll() ([]dto.RoleResponse, global.ErrorResponse)
	FindById(id string) (*dto.RoleResponse, global.ErrorResponse)
}

type roleUsecase struct {
	roleRepository repository.RoleRepository
}

func NewRoleUsecase(roleRepository repository.RoleRepository) RoleUsecase {
	return &roleUsecase{roleRepository}
}

func (u *roleUsecase) FindAll() ([]dto.RoleResponse, global.ErrorResponse) {
	roles, err := u.roleRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return mapper.ToRoleResponses(roles), nil
}

func (u *roleUsecase) FindById(id string) (*dto.RoleResponse, global.ErrorResponse) {
	role, err := u.roleRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToRoleResponse(role), nil
}

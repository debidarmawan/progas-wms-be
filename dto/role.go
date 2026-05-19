package dto

type GetRoleByIdSpec struct {
	Id string `json:"id"`
}

type RoleResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

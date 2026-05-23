package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=6"`
	RoleId   string `json:"role_id" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	RoleId   string `json:"role_id" validate:"required"`
	Password string `json:"password" validate:"omitempty,min=6"`
	IsActive bool   `json:"is_active"`
}

type UserListResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	RoleId         string `json:"role_id"`
	RoleName       string `json:"role_name"`
	IsActive       bool   `json:"is_active"`
	LastLoggedInAt string `json:"last_logged_in_at,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type PaginatedUserList struct {
	Items []UserListResponse `json:"items"`
	Meta  PaginationMeta     `json:"meta"`
}

package handler

import (
	"progas-wms-be/global"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type RoleHandler struct {
	roleUsecase usecase.RoleUsecase
}

func NewRoleHandler(roleUsecase usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUsecase}
}

func (h *RoleHandler) Routes(group fiber.Router) {
	group.Get("/roles/", h.FindAll)
	group.Get("/roles/:id", h.FindById)
}

// FindAll godoc
//
//	@Summary		Find all roles
//	@Description	Find all roles
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}	global.Response[[]dto.RoleResponse]
//	@Router			/roles [get]
func (h *RoleHandler) FindAll(c fiber.Ctx) error {
	res, err := h.roleUsecase.FindAll()
	if err != nil {
		return err.ToResponse(c)
	}

	return global.CreateResponse(&res, fiber.StatusOK, c)
}

// FindById godoc
//
//	@Summary		Find role by id
//	@Description	Find role by id
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Role ID"
//	@Success		200	{object}	global.Response[dto.RoleResponse]
//	@Router			/roles/:id [get]
func (h *RoleHandler) FindById(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.roleUsecase.FindById(id)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(&res, fiber.StatusOK, c)
}

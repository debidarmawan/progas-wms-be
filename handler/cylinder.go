package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type CylinderHandler struct {
	usecase usecase.CylinderUsecase
}

func NewCylinderHandler(usecase usecase.CylinderUsecase) *CylinderHandler {
	return &CylinderHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all cylinders
//	@Description	Find all registered cylinders with pagination and search
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by barcode, status, ownership, or item name/SKU/gas type"
//	@Success		200		{object}	global.Response[dto.PaginatedCylinderList]
//	@Router			/cylinders [get]
func (h *CylinderHandler) FindAll(c fiber.Ctx) error {
	var query dto.ListQuery
	if err := helper.ValidateQuery(c, &query); err != nil {
		return err.ToResponse(c)
	}

	res, err := h.usecase.FindAll(&query)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// FindById godoc
//
//	@Summary		Find cylinder by id
//	@Description	Find cylinder by id
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Cylinder ID"
//	@Success		200	{object}	global.Response[dto.CylinderResponse]
//	@Router			/cylinders/{id} [get]
func (h *CylinderHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Register cylinder
//	@Description	Register a new cylinder with unique barcode
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateCylinderRequest	true	"Register cylinder request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/cylinders [post]
func (h *CylinderHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCylinderRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Cylinder registered successfully", fiber.StatusOK, c)
}

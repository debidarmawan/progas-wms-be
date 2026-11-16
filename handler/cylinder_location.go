package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type CylinderLocationHandler struct {
	usecase usecase.CylinderLocationUsecase
}

func NewCylinderLocationHandler(usecase usecase.CylinderLocationUsecase) *CylinderLocationHandler {
	return &CylinderLocationHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all cylinder locations
//	@Description	Find all company-owned cylinder locations (warehouse, filling station, workshop)
//	@Tags			Cylinder Location
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by name or code"
//	@Success		200		{object}	global.Response[dto.PaginatedCylinderLocationList]
//	@Router			/cylinder-locations [get]
func (h *CylinderLocationHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Find cylinder location by id
//	@Description	Find a cylinder location by id
//	@Tags			Cylinder Location
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Cylinder location ID"
//	@Success		200	{object}	global.Response[dto.CylinderLocationResponse]
//	@Router			/cylinder-locations/{id} [get]
func (h *CylinderLocationHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create cylinder location
//	@Description	Create a cylinder location (warehouse, filling station or workshop)
//	@Tags			Cylinder Location
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateCylinderLocationRequest	true	"Create cylinder location request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/cylinder-locations [post]
func (h *CylinderLocationHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCylinderLocationRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Cylinder location created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update cylinder location
//	@Description	Update a cylinder location by id
//	@Tags			Cylinder Location
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string								true	"Cylinder location ID"
//	@Param			request	body		dto.UpdateCylinderLocationRequest	true	"Update cylinder location request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/cylinder-locations/{id} [put]
func (h *CylinderLocationHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateCylinderLocationRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Cylinder location updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete cylinder location
//	@Description	Delete a cylinder location by id (blocked while still assigned to cylinders)
//	@Tags			Cylinder Location
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Cylinder location ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/cylinder-locations/{id} [delete]
func (h *CylinderLocationHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Cylinder location deleted successfully", fiber.StatusOK, c)
}

package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type UomHandler struct {
	usecase usecase.UomUsecase
}

func NewUomHandler(usecase usecase.UomUsecase) *UomHandler {
	return &UomHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all UOM
//	@Description	Find all units of measure with pagination and search
//	@Tags			UOM
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by name or code"
//	@Success		200		{object}	global.Response[dto.PaginatedUomList]
//	@Router			/uoms [get]
func (h *UomHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Find UOM by id
//	@Description	Find unit of measure by id
//	@Tags			UOM
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"UOM ID"
//	@Success		200	{object}	global.Response[dto.UomResponse]
//	@Router			/uoms/{id} [get]
func (h *UomHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create UOM
//	@Description	Create a unit of measure
//	@Tags			UOM
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateUomRequest	true	"Create UOM request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/uoms [post]
func (h *UomHandler) Create(c fiber.Ctx) error {
	var req dto.CreateUomRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("UOM created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update UOM
//	@Description	Update a unit of measure by id
//	@Tags			UOM
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string					true	"UOM ID"
//	@Param			request	body		dto.UpdateUomRequest	true	"Update UOM request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/uoms/{id} [put]
func (h *UomHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateUomRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("UOM updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete UOM
//	@Description	Delete a unit of measure by id (blocked when still used by master items)
//	@Tags			UOM
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"UOM ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/uoms/{id} [delete]
func (h *UomHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("UOM deleted successfully", fiber.StatusOK, c)
}

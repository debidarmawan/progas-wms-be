package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type MasterItemHandler struct {
	usecase usecase.MasterItemUsecase
}

func NewMasterItemHandler(usecase usecase.MasterItemUsecase) *MasterItemHandler {
	return &MasterItemHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all master items
//	@Description	Find all master items with pagination and search
//	@Tags			Master Item
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by name, SKU, or gas type"
//	@Success		200		{object}	global.Response[dto.PaginatedMasterItemList]
//	@Router			/master-items [get]
func (h *MasterItemHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Find master item by id
//	@Description	Find master item by id
//	@Tags			Master Item
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Master Item ID"
//	@Success		200	{object}	global.Response[dto.MasterItemResponse]
//	@Router			/master-items/{id} [get]
func (h *MasterItemHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create master item
//	@Description	Create master item (serialized gas or non-serialized spare part)
//	@Tags			Master Item
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateMasterItemRequest	true	"Create master item request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/master-items [post]
func (h *MasterItemHandler) Create(c fiber.Ctx) error {
	var req dto.CreateMasterItemRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Master item created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update master item
//	@Description	Update master item by id
//	@Tags			Master Item
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string						true	"Master Item ID"
//	@Param			request	body		dto.UpdateMasterItemRequest	true	"Update master item request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/master-items/{id} [put]
func (h *MasterItemHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateMasterItemRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Master item updated successfully", fiber.StatusOK, c)
}

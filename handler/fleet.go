package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type FleetHandler struct {
	usecase usecase.FleetUsecase
}

func NewFleetHandler(usecase usecase.FleetUsecase) *FleetHandler {
	return &FleetHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		List fleet vehicles
//	@Description	List fleet vehicles with pagination and search
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by plate number or driver name"
//	@Success		200		{object}	global.Response[dto.PaginatedFleetList]
//	@Router			/logistics/fleet [get]
func (h *FleetHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Get fleet vehicle by id
//	@Description	Get fleet vehicle detail
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Fleet vehicle ID"
//	@Success		200	{object}	global.Response[dto.FleetResponse]
//	@Router			/logistics/fleet/{id} [get]
func (h *FleetHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create fleet vehicle
//	@Description	Register a fleet vehicle with max weight capacity
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateFleetRequest	true	"Create fleet request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/logistics/fleet [post]
func (h *FleetHandler) Create(c fiber.Ctx) error {
	var req dto.CreateFleetRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Fleet vehicle created", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update fleet vehicle
//	@Description	Update fleet vehicle driver, max weight, or active status
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string					true	"Fleet vehicle ID"
//	@Param			request	body		dto.UpdateFleetRequest	true	"Update fleet request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/logistics/fleet/{id} [put]
func (h *FleetHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateFleetRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Fleet vehicle updated", fiber.StatusOK, c)
}

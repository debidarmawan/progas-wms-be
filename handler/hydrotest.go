package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type HydrotestHandler struct {
	usecase usecase.HydrotestUsecase
}

func NewHydrotestHandler(usecase usecase.HydrotestUsecase) *HydrotestHandler {
	return &HydrotestHandler{usecase: usecase}
}

// FindDue godoc
//
//	@Summary	List cylinders due for hydrotest
//	@Tags		Maintenance
//	@Security	Bearer
//	@Param		days	query		int	false	"Due within days (default 30)"
//	@Success	200		{object}	global.Response[dto.HydrotestDueResponse]
//	@Router		/maintenance/hydrotest/due [get]
func (h *HydrotestHandler) FindDue(c fiber.Ctx) error {
	days := 30
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil {
			days = v
		}
	}
	res, err := h.usecase.FindDue(days)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Record godoc
//
//	@Summary		Record cylinder hydrotest
//	@Description	Update last hydrotest date and set status to MAINTENANCE
//	@Tags			Maintenance
//	@Security		Bearer
//	@Param			id		path		string						true	"Cylinder ID"
//	@Param			request	body		dto.RecordHydrotestRequest	true	"Request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/maintenance/cylinders/{id}/hydrotest [post]
func (h *HydrotestHandler) Record(c fiber.Ctx) error {
	var req dto.RecordHydrotestRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Record(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Hydrotest recorded", fiber.StatusOK, c)
}

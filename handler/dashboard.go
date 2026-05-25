package handler

import (
	"progas-wms-be/global"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type DashboardHandler struct {
	usecase usecase.DashboardUsecase
}

func NewDashboardHandler(usecase usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{usecase: usecase}
}

// GetSummary godoc
//
//	@Summary		Dashboard summary
//	@Description	Stock, outstanding, low stock alerts, hydrotest alerts, quota alerts
//	@Tags			Dashboard
//	@Security		Bearer
//	@Success		200	{object}	global.Response[dto.DashboardSummaryResponse]
//	@Router			/dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c fiber.Ctx) error {
	res, err := h.usecase.GetSummary()
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

package handler

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/repository"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type ExchangeHandler struct {
	usecase  usecase.ExchangeUsecase
	rbacRepo repository.RbacRepository
}

func NewExchangeHandler(usecase usecase.ExchangeUsecase, rbacRepo repository.RbacRepository) *ExchangeHandler {
	return &ExchangeHandler{usecase: usecase, rbacRepo: rbacRepo}
}

func (h *ExchangeHandler) canApprove(c fiber.Ctx) bool {
	roleId, _ := c.Locals("role_id").(string)
	if roleId == "" {
		return false
	}
	isSuperAdmin, err := h.rbacRepo.IsSuperAdmin(roleId)
	if err == nil && isSuperAdmin {
		return true
	}
	allowed, err := h.rbacRepo.HasPermission(roleId, constant.PermExchangeApprove)
	return err == nil && allowed
}

// Process godoc
//
//	@Summary		Process cylinder exchange
//	@Description	Gate swap: OUT (IN_TRANSIT→OUTSTANDING) and IN (OUTSTANDING→EMPTY). Updates outstanding count (excludes CUSTOMER-owned cylinders). Use force_approve with exchange.approve permission when over quota.
//	@Tags			Outbound
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.ProcessExchangeRequest	true	"Exchange request"
//	@Success		200		{object}	global.Response[dto.ExchangeResponse]
//	@Router			/outbound/exchange [post]
func (h *ExchangeHandler) Process(c fiber.Ctx) error {
	var req dto.ProcessExchangeRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Process(actorUserId, &req, h.canApprove(c))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

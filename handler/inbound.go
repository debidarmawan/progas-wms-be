package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type InboundHandler struct {
	usecase usecase.InboundUsecase
}

func NewInboundHandler(usecase usecase.InboundUsecase) *InboundHandler {
	return &InboundHandler{usecase: usecase}
}

// EmptyReceive godoc
//
//	@Summary		Empty cylinder receiving
//	@Description	Receive empty cylinders returned from customer (OUTSTANDING/IN_TRANSIT → EMPTY)
//	@Tags			Inbound
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.BarcodeListRequest	true	"Barcode list"
//	@Success		200		{object}	global.Response[dto.BarcodeOperationResponse]
//	@Router			/inbound/empty-receive [post]
func (h *InboundHandler) EmptyReceive(c fiber.Ctx) error {
	var req dto.BarcodeListRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.EmptyReceive(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// PreFillQC godoc
//
//	@Summary		Pre-fill QC
//	@Description	Pass pre-fill inspection (EMPTY → READY_TO_FILL)
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.BarcodeListRequest	true	"Barcode list"
//	@Success		200		{object}	global.Response[dto.BarcodeOperationResponse]
//	@Router			/production/qc/pre-fill [post]
func (h *InboundHandler) PreFillQC(c fiber.Ctx) error {
	var req dto.BarcodeListRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.PreFillQC(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// PostFillQC godoc
//
//	@Summary		Post-fill QC
//	@Description	Pass post-fill inspection after gas filling (FILLED → READY)
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.BarcodeListRequest	true	"Barcode list"
//	@Success		200		{object}	global.Response[dto.BarcodeOperationResponse]
//	@Router			/production/qc/post-fill [post]
func (h *InboundHandler) PostFillQC(c fiber.Ctx) error {
	var req dto.BarcodeListRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.PostFillQC(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

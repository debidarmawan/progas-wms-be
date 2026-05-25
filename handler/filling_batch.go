package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type FillingBatchHandler struct {
	usecase usecase.FillingBatchUsecase
}

func NewFillingBatchHandler(usecase usecase.FillingBatchUsecase) *FillingBatchHandler {
	return &FillingBatchHandler{usecase: usecase}
}

// Submit godoc
//
//	@Summary		Submit filling batch
//	@Description	Create and complete a filling batch atomically (validates status & cross-gas, sets cylinders to FILLED; use post-fill QC for READY)
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.SubmitFillingBatchRequest	true	"Filling batch request"
//	@Success		200		{object}	global.Response[dto.FillingBatchResponse]
//	@Router			/production/filling-batches [post]
func (h *FillingBatchHandler) Submit(c fiber.Ctx) error {
	var req dto.SubmitFillingBatchRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Submit(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// FindAll godoc
//
//	@Summary		List filling batches
//	@Description	List filling batches with pagination and search
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by batch number or gas type"
//	@Success		200		{object}	global.Response[dto.PaginatedFillingBatchList]
//	@Router			/production/filling-batches [get]
func (h *FillingBatchHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Get filling batch by id
//	@Description	Get filling batch detail including scanned cylinders
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Filling batch ID"
//	@Success		200	{object}	global.Response[dto.FillingBatchResponse]
//	@Router			/production/filling-batches/{id} [get]
func (h *FillingBatchHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

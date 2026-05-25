package handler

import (
	"progas-wms-be/global"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type ReportHandler struct {
	usecase usecase.ReportUsecase
}

func NewReportHandler(usecase usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{usecase: usecase}
}

// StockLedger godoc
//
//	@Summary		Stock ledger per barcode
//	@Description	Status change history for a cylinder barcode
//	@Tags			Report
//	@Security		Bearer
//	@Param			barcode	query		string	true	"Cylinder barcode SN"
//	@Success		200		{object}	global.Response[dto.StockLedgerReportResponse]
//	@Router			/reports/stock-ledger [get]
func (h *ReportHandler) StockLedger(c fiber.Ctx) error {
	res, err := h.usecase.StockLedger(c.Query("barcode"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Turnaround godoc
//
//	@Summary		Turn-around rate report
//	@Description	Average days from EMPTY to READY per cylinder in date range
//	@Tags			Report
//	@Security		Bearer
//	@Param			from	query		string	false	"From date YYYY-MM-DD"
//	@Param			to		query		string	false	"To date YYYY-MM-DD"
//	@Success		200		{object}	global.Response[dto.TurnaroundReportResponse]
//	@Router			/reports/turnaround [get]
func (h *ReportHandler) Turnaround(c fiber.Ctx) error {
	res, err := h.usecase.Turnaround(c.Query("from"), c.Query("to"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

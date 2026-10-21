package handler

import (
	"bytes"
	"encoding/csv"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type CylinderHandler struct {
	usecase usecase.CylinderUsecase
}

func NewCylinderHandler(usecase usecase.CylinderUsecase) *CylinderHandler {
	return &CylinderHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all cylinders
//	@Description	Find all registered cylinders with pagination and search
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by barcode, status, ownership, or item name/SKU/gas type"
//	@Success		200		{object}	global.Response[dto.PaginatedCylinderList]
//	@Router			/cylinders [get]
func (h *CylinderHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Find cylinder by id
//	@Description	Find cylinder by id
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Cylinder ID"
//	@Success		200	{object}	global.Response[dto.CylinderResponse]
//	@Router			/cylinders/{id} [get]
func (h *CylinderHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// FindByBarcode godoc
//
//	@Summary		Find cylinder by barcode
//	@Description	Lookup a single cylinder by its barcode serial number (for scan validation)
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			sn	path		string	true	"Cylinder barcode serial number"
//	@Success		200	{object}	global.Response[dto.CylinderResponse]
//	@Router			/cylinders/by-barcode/{sn} [get]
func (h *CylinderHandler) FindByBarcode(c fiber.Ctx) error {
	res, err := h.usecase.FindByBarcode(c.Params("sn"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// ResolveScan godoc
//
//	@Summary		Resolve a scanned barcode for a menu
//	@Description	Return the cylinder plus whether its status allows the scanned menu (single source of truth for mobile)
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			sn		path		string	true	"Cylinder barcode serial number"
//	@Param			menu	query		string	true	"pre-fill | filling | post-fill | do | exchange-out | exchange-in | hydrotest"
//	@Success		200		{object}	global.Response[dto.ScanResolveResponse]
//	@Router			/cylinders/scan/{sn} [get]
func (h *CylinderHandler) ResolveScan(c fiber.Ctx) error {
	var query dto.ScanCylinderQuery
	if err := helper.ValidateQuery(c, &query); err != nil {
		return err.ToResponse(c)
	}
	res, err := h.usecase.ResolveScan(c.Params("sn"), query.Menu)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Labels godoc
//
//	@Summary		Export cylinder labels
//	@Description	CSV of barcode serials to print. Pass ids=... or labeled=false for every never-printed label.
//	@Tags			Cylinder
//	@Produce		text/csv
//	@Security		Bearer
//	@Param			ids		query		string	false	"Comma-separated cylinder IDs"
//	@Param			labeled	query		string	false	"Set false to list cylinders whose label was never printed"
//	@Param			limit	query		int		false	"Max rows when listing unlabeled (default 500)"
//	@Param			format	query		string	false	"json to return rows as JSON (default csv)"
//	@Success		200		{string}	string	"CSV: barcode_sn,item_name,gas_type (or JSON envelope when format=json)"
//	@Router			/cylinders/labels [get]
func (h *CylinderHandler) Labels(c fiber.Ctx) error {
	var query dto.CylinderLabelQuery
	if err := helper.ValidateQuery(c, &query); err != nil {
		return err.ToResponse(c)
	}

	rows, err := h.usecase.FindLabels(strings.Split(query.Ids, ","), query.Labeled == "false", query.Limit)
	if err != nil {
		return err.ToResponse(c)
	}

	if query.Format == "json" {
		return global.CreateResponse(&rows, fiber.StatusOK, c)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	_ = writer.Write([]string{"barcode_sn", "item_name", "gas_type"})
	for _, row := range rows {
		_ = writer.Write([]string{row.BarcodeSN, row.ItemName, row.GasType})
	}
	writer.Flush()
	if writeErr := writer.Error(); writeErr != nil {
		return global.InternalServerError(writeErr).ToResponse(c)
	}

	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="cylinder-labels.csv"`)
	return c.SendString(buf.String())
}

// MarkLabelsPrinted godoc
//
//	@Summary		Mark labels as printed
//	@Description	Record the print time so barcode_sn can no longer be edited
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.MarkLabelsPrintedRequest	true	"Printed cylinder IDs"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/cylinders/labels/printed [post]
func (h *CylinderHandler) MarkLabelsPrinted(c fiber.Ctx) error {
	var req dto.MarkLabelsPrintedRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.MarkLabelsPrinted(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Label tercatat tercetak", fiber.StatusOK, c)
}

// History godoc
//
//	@Summary		Get cylinder change history
//	@Description	Get history of master data changes (product, ownership, hydrotest date, remarks) and status changes for a cylinder
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Cylinder ID"
//	@Success		200	{object}	global.Response[dto.CylinderHistoryResponse]
//	@Router			/cylinders/{id}/history [get]
func (h *CylinderHandler) History(c fiber.Ctx) error {
	res, err := h.usecase.History(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Register cylinder
//	@Description	Register a new cylinder with unique barcode
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateCylinderRequest	true	"Register cylinder request"
//	@Success		200		{object}	global.Response[dto.CylinderResponse]
//	@Router			/cylinders [post]
func (h *CylinderHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCylinderRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Create(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update cylinder
//	@Description	Update an existing cylinder's metadata
//	@Tags			Cylinder
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string						true	"Cylinder ID"
//	@Param			request	body		dto.UpdateCylinderRequest	true	"Update cylinder request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/cylinders/{id} [put]
func (h *CylinderHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateCylinderRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Tabung diperbarui", fiber.StatusOK, c)
}

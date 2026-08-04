package mapper

import (
	"encoding/json"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

var cylinderHistoryActionLabels = map[string]string{
	constant.AuditCylinderCreate:      "Registrasi",
	constant.AuditCylinderUpdate:      "Perubahan Data",
	constant.LedgerActionEmptyReceive: "Terima Kosong",
	constant.LedgerActionPreFillQC:    "QC Pre-Fill",
	constant.LedgerActionFillingBatch: "Filling Batch",
	constant.LedgerActionPostFillQC:   "QC Post-Fill",
	constant.LedgerActionDOIssue:      "Terbit Surat Jalan",
	constant.LedgerActionExchangeOut:  "Tukar Keluar",
	constant.LedgerActionExchangeIn:   "Tukar Masuk",
	constant.LedgerActionHydrotest:    "Hydrotest",
}

func CylinderHistoryActionLabel(action string) string {
	if label, ok := cylinderHistoryActionLabels[action]; ok {
		return label
	}
	return action
}

func ParseCylinderHistoryChanges(details string, itemName func(string) string) []dto.CylinderHistoryChange {
	if details == "" {
		return []dto.CylinderHistoryChange{}
	}
	var wrapper struct {
		Changes []dto.CylinderHistoryChange `json:"changes"`
	}
	if err := json.Unmarshal([]byte(details), &wrapper); err == nil && wrapper.Changes != nil {
		return wrapper.Changes
	}

	var legacy struct {
		BarcodeSN     string `json:"barcode_sn"`
		OwnershipType string `json:"ownership_type"`
		ItemId        string `json:"item_id"`
	}
	if err := json.Unmarshal([]byte(details), &legacy); err != nil {
		return []dto.CylinderHistoryChange{}
	}

	changes := []dto.CylinderHistoryChange{}
	if legacy.BarcodeSN != "" {
		changes = append(changes, dto.CylinderHistoryChange{Field: "barcode_sn", Label: "Barcode", New: legacy.BarcodeSN})
	}
	if legacy.ItemId != "" {
		name := ""
		if itemName != nil {
			name = itemName(legacy.ItemId)
		}
		changes = append(changes, dto.CylinderHistoryChange{Field: "item_id", Label: "Produk Gas", New: name})
	}
	if legacy.OwnershipType != "" {
		changes = append(changes, dto.CylinderHistoryChange{Field: "ownership_type", Label: "Kepemilikan", New: legacy.OwnershipType})
	}
	return changes
}

func ToCylinderHistoryResponse(cylinder *model.Cylinder, entries []dto.CylinderHistoryEntry) *dto.CylinderHistoryResponse {
	if entries == nil {
		entries = []dto.CylinderHistoryEntry{}
	}
	return &dto.CylinderHistoryResponse{
		CylinderId: cylinder.Id,
		BarcodeSN:  cylinder.BarcodeSN,
		Entries:    entries,
	}
}

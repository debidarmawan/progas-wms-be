package helper

import (
	"fmt"
	"progas-wms-be/enum"
	"progas-wms-be/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ValidateBarcodeList(barcodes []string) error {
	seen := make(map[string]bool)
	for _, b := range barcodes {
		b = strings.TrimSpace(b)
		if b == "" {
			return fmt.Errorf("barcode cannot be empty")
		}
		if seen[b] {
			return fmt.Errorf("duplicate barcode: %s", b)
		}
		seen[b] = true
	}
	if len(seen) == 0 {
		return fmt.Errorf("barcodes are required")
	}
	return nil
}

func GenerateFillingBatchNumber() string {
	id, _ := uuid.NewV7()
	return fmt.Sprintf("FB-%s-%s", time.Now().Format("20060102"), id.String()[:8])
}

func CanReceiveAsEmpty(status enum.CylinderStatus) bool {
	return status == enum.CylinderStatusOutstanding || status == enum.CylinderStatusInTransit
}

func CanPreFillQC(status enum.CylinderStatus) bool {
	return status == enum.CylinderStatusEmpty
}

func CanEnterFillingBatch(status enum.CylinderStatus) bool {
	return status == enum.CylinderStatusEmpty || status == enum.CylinderStatusReadyToFill
}

func ValidateFillingBatchCylinders(cylinders []model.Cylinder, batchItem *model.MasterItem) error {
	if !batchItem.IsSerialized {
		return fmt.Errorf("batch item must be a serialized gas product")
	}
	if batchItem.GasType == "" {
		return fmt.Errorf("batch item has no gas type configured")
	}

	for _, cyl := range cylinders {
		if !CanEnterFillingBatch(cyl.Status) {
			return fmt.Errorf("cylinder %s has invalid status %s for filling (must be EMPTY or READY_TO_FILL)", cyl.BarcodeSN, cyl.Status)
		}
		if cyl.MasterItem.GasType != batchItem.GasType {
			return fmt.Errorf("cylinder %s gas type %s does not match batch gas type %s", cyl.BarcodeSN, cyl.MasterItem.GasType, batchItem.GasType)
		}
	}
	return nil
}

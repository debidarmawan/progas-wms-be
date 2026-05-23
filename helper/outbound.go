package helper

import (
	"fmt"
	"progas-wms-be/enum"
	"progas-wms-be/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateDONumber() string {
	id, _ := uuid.NewV7()
	return fmt.Sprintf("DO-%s-%s", time.Now().Format("20060102"), id.String()[:8])
}

func GenerateWorkOrderNumber() string {
	id, _ := uuid.NewV7()
	return fmt.Sprintf("WO-%s-%s", time.Now().Format("20060102"), id.String()[:8])
}

func CylinderFilledWeightKg(item model.MasterItem) float64 {
	return item.EmptyWeightKg + item.GasWeightKg
}

func SumCylinderWeight(cylinders []model.Cylinder) float64 {
	var total float64
	for _, cyl := range cylinders {
		total += CylinderFilledWeightKg(cyl.MasterItem)
	}
	return total
}

func CanIssueOnDO(status enum.CylinderStatus) bool {
	return status == enum.CylinderStatusReady
}

func CanDeliverOnExchange(status enum.CylinderStatus) bool {
	return status == enum.CylinderStatusInTransit
}

func CanReceiveOnExchange(status enum.CylinderStatus) bool {
	return status == enum.CylinderStatusOutstanding
}

func CountsTowardOutstanding(cyl model.Cylinder) bool {
	return cyl.OwnershipType != enum.OwnershipCustomer
}

func CountOutstandingDelta(cylinders []model.Cylinder) int {
	count := 0
	for _, cyl := range cylinders {
		if CountsTowardOutstanding(cyl) {
			count++
		}
	}
	return count
}

func ValidateDOCylinders(cylinders []model.Cylinder) error {
	for _, cyl := range cylinders {
		if !CanIssueOnDO(cyl.Status) {
			return fmt.Errorf("cylinder %s has invalid status %s for DO (must be READY)", cyl.BarcodeSN, cyl.Status)
		}
		if !cyl.MasterItem.IsSerialized {
			return fmt.Errorf("cylinder %s is not a serialized gas product", cyl.BarcodeSN)
		}
	}
	return nil
}

func ValidateExchangeOutCylinders(cylinders []model.Cylinder) error {
	for _, cyl := range cylinders {
		if !CanDeliverOnExchange(cyl.Status) {
			return fmt.Errorf("cylinder %s has invalid status %s for exchange OUT (must be IN_TRANSIT)", cyl.BarcodeSN, cyl.Status)
		}
	}
	return nil
}

func ValidateExchangeInCylinders(cylinders []model.Cylinder) error {
	for _, cyl := range cylinders {
		if !CanReceiveOnExchange(cyl.Status) {
			return fmt.Errorf("cylinder %s has invalid status %s for exchange IN (must be OUTSTANDING)", cyl.BarcodeSN, cyl.Status)
		}
	}
	return nil
}

func DetectCrossCustomerAlerts(customerId string, cylinders []model.Cylinder) []string {
	var alerts []string
	for _, cyl := range cylinders {
		if cyl.OwnerId != nil && *cyl.OwnerId != "" && *cyl.OwnerId != customerId {
			alerts = append(alerts, fmt.Sprintf("cylinder %s registered to different customer (owner_id=%s)", cyl.BarcodeSN, *cyl.OwnerId))
		}
	}
	return alerts
}

func ValidateBarcodeListsUnique(outBarcodes, inBarcodes []string) error {
	if err := ValidateBarcodeList(outBarcodes); err != nil {
		return fmt.Errorf("out_barcodes: %w", err)
	}
	if err := ValidateBarcodeList(inBarcodes); err != nil {
		return fmt.Errorf("in_barcodes: %w", err)
	}
	seen := make(map[string]bool)
	for _, b := range outBarcodes {
		b = strings.TrimSpace(b)
		if seen[b] {
			return fmt.Errorf("barcode appears in both out and in lists: %s", b)
		}
		seen[b] = true
	}
	for _, b := range inBarcodes {
		b = strings.TrimSpace(b)
		if seen[b] {
			return fmt.Errorf("barcode appears in both out and in lists: %s", b)
		}
	}
	return nil
}

func WouldExceedQuota(currentOutstanding, quotaLimit, outDelta int) bool {
	if quotaLimit <= 0 {
		return false
	}
	return currentOutstanding+outDelta > quotaLimit
}

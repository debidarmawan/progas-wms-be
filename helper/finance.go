package helper

import (
	"fmt"
	"progas-wms-be/enum"
	"time"

	"github.com/google/uuid"
)

func GenerateInvoiceNumber() string {
	id, _ := uuid.NewV7()
	return fmt.Sprintf("INV-%s-%s", time.Now().Format("20060102"), id.String()[:8])
}

// DetermineInvoiceStatus decides UNPAID/PARTIAL/PAID for an invoice given its
// total amount and the paid amount after applying a new payment.
func DetermineInvoiceStatus(paidAmount, totalAmount float64) enum.InvoiceStatus {
	if paidAmount <= 0 {
		return enum.InvoiceStatusUnpaid
	}
	if paidAmount >= totalAmount {
		return enum.InvoiceStatusPaid
	}
	return enum.InvoiceStatusPartial
}

package helper

import (
	"progas-wms-be/enum"
	"testing"
)

func TestDetermineInvoiceStatus(t *testing.T) {
	cases := []struct {
		name       string
		paidAmount float64
		total      float64
		want       enum.InvoiceStatus
	}{
		{"no payment yet", 0, 100000, enum.InvoiceStatusUnpaid},
		{"partial payment", 40000, 100000, enum.InvoiceStatusPartial},
		{"exact full payment", 100000, 100000, enum.InvoiceStatusPaid},
		{"overpayment still counts as paid", 100001, 100000, enum.InvoiceStatusPaid},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DetermineInvoiceStatus(c.paidAmount, c.total)
			if got != c.want {
				t.Errorf("DetermineInvoiceStatus(%v, %v) = %v, want %v", c.paidAmount, c.total, got, c.want)
			}
		})
	}
}

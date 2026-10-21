package helper

import (
	"progas-wms-be/enum"
	"reflect"
	"testing"
)

func TestCanScanMenu(t *testing.T) {
	cases := []struct {
		name   string
		menu   string
		status enum.CylinderStatus
		want   bool
	}{
		{"pre-fill accepts EMPTY", "pre-fill", enum.CylinderStatusEmpty, true},
		{"pre-fill rejects READY_TO_FILL", "pre-fill", enum.CylinderStatusReadyToFill, false},
		{"filling accepts READY_TO_FILL", "filling", enum.CylinderStatusReadyToFill, true},
		{"filling rejects READY", "filling", enum.CylinderStatusReady, false},
		{"post-fill accepts FILLED", "post-fill", enum.CylinderStatusFilled, true},
		{"post-fill rejects EMPTY", "post-fill", enum.CylinderStatusEmpty, false},
		{"do accepts READY", "do", enum.CylinderStatusReady, true},
		{"do rejects FILLED", "do", enum.CylinderStatusFilled, false},
		{"exchange-out accepts IN_TRANSIT", "exchange-out", enum.CylinderStatusInTransit, true},
		{"exchange-out rejects READY", "exchange-out", enum.CylinderStatusReady, false},
		{"exchange-in accepts OUTSTANDING", "exchange-in", enum.CylinderStatusOutstanding, true},
		{"exchange-in rejects EMPTY", "exchange-in", enum.CylinderStatusEmpty, false},
		{"hydrotest accepts EMPTY", "hydrotest", enum.CylinderStatusEmpty, true},
		{"hydrotest rejects LOST", "hydrotest", enum.CylinderStatusLost, false},
		{"hydrotest rejects WRITE_OFF", "hydrotest", enum.CylinderStatusWriteOff, false},
		{"unknown menu always rejected", "stock-adjust", enum.CylinderStatusEmpty, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CanScanMenu(c.menu, c.status); got != c.want {
				t.Errorf("CanScanMenu(%q, %v) = %v, want %v", c.menu, c.status, got, c.want)
			}
		})
	}
}

func TestIsScanMenu(t *testing.T) {
	cases := []struct {
		name string
		menu string
		want bool
	}{
		{"pre-fill", "pre-fill", true},
		{"exchange-in", "exchange-in", true},
		{"case sensitive", "PRE-FILL", false},
		{"unknown", "trip", false},
		{"empty", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsScanMenu(c.menu); got != c.want {
				t.Errorf("IsScanMenu(%q) = %v, want %v", c.menu, got, c.want)
			}
		})
	}
}

func TestAllowedScanMenus(t *testing.T) {
	cases := []struct {
		name   string
		status enum.CylinderStatus
		want   []string
	}{
		{"EMPTY can be pre-filled, filled, or hydrotested", enum.CylinderStatusEmpty, []string{"pre-fill", "filling", "hydrotest"}},
		{"READY can go on DO", enum.CylinderStatusReady, []string{"do", "hydrotest"}},
		{"OUTSTANDING can only come back in", enum.CylinderStatusOutstanding, []string{"exchange-in", "hydrotest"}},
		{"LOST cannot be scanned anywhere", enum.CylinderStatusLost, []string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := AllowedScanMenus(c.status); !reflect.DeepEqual(got, c.want) {
				t.Errorf("AllowedScanMenus(%v) = %v, want %v", c.status, got, c.want)
			}
		})
	}
}

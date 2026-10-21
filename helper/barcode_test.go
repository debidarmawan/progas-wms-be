package helper

import (
	"strings"
	"testing"
)

func TestNormalizeBarcode(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"trims whitespace", "  abc-123  ", "ABC-123"},
		{"uppercases scanner output", "pg-co2-7k2m9q", "PG-CO2-7K2M9Q"},
		{"empty stays empty", "   ", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NormalizeBarcode(c.in); got != c.want {
				t.Errorf("NormalizeBarcode(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestIsValidBarcode(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"generated serial", "PG-CO2-7K2M9Q", true},
		{"numeric legacy serial", "12345678", true},
		{"too short", "AB-1", false},
		{"lowercase rejected", "pg-co2-7k2m9q", false},
		{"space rejected", "PG CO2 123", false},
		{"underscore rejected", "PG_CO2_123", false},
		{"over 32 chars", strings.Repeat("A", 33), false},
		{"empty rejected", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsValidBarcode(c.in); got != c.want {
				t.Errorf("IsValidBarcode(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestGenerateCylinderBarcode(t *testing.T) {
	got := GenerateCylinderBarcode("CO2")
	if !IsValidBarcode(got) {
		t.Fatalf("generated barcode %q is not valid", got)
	}
	if !strings.HasPrefix(got, "PG-CO2-") {
		t.Errorf("GenerateCylinderBarcode() = %q, want prefix PG-CO2-", got)
	}

	suffix := got[len(got)-6:]
	for _, ambiguous := range []string{"I", "L", "O", "U"} {
		if strings.Contains(suffix, ambiguous) {
			t.Errorf("generated suffix %q contains ambiguous char %q", suffix, ambiguous)
		}
	}
	if a, b := GenerateCylinderBarcode("CO2"), GenerateCylinderBarcode("CO2"); a == b {
		t.Errorf("two generated barcodes collided: %q", a)
	}
	if got := GenerateCylinderBarcode(""); !strings.HasPrefix(got, "PG-CYL-") {
		t.Errorf("GenerateCylinderBarcode(\"\") = %q, want PG-CYL- prefix", got)
	}
}

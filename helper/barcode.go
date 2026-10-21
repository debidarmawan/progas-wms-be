package helper

import (
	"crypto/rand"
	"regexp"
	"strings"
)

// Strict format for NEWLY registered serials. Legacy serials already stamped in
// the field are only normalized, never rejected on lookup.
var barcodePattern = regexp.MustCompile(`^[A-Z0-9-]{6,32}$`)

// Crockford base32: no I, L, O, U so the printed label stays readable by hand.
const barcodeAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// NormalizeBarcode makes scanned and typed input comparable: scanners and humans
// both add whitespace and case noise.
func NormalizeBarcode(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func IsValidBarcode(s string) bool {
	return barcodePattern.MatchString(s)
}

// GenerateCylinderBarcode returns PG-<SKU4>-<6 random chars>; caller retries on
// unique-index conflict instead of keeping a counter table.
func GenerateCylinderBarcode(sku string) string {
	var prefix strings.Builder
	for _, r := range strings.ToUpper(sku) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			prefix.WriteRune(r)
		}
		if prefix.Len() == 4 {
			break
		}
	}
	if prefix.Len() == 0 {
		prefix.WriteString("CYL")
	}

	random := make([]byte, 6)
	if _, err := rand.Read(random); err != nil {
		return ""
	}
	suffix := make([]byte, len(random))
	for i, v := range random {
		suffix[i] = barcodeAlphabet[int(v)%len(barcodeAlphabet)]
	}
	return "PG-" + prefix.String() + "-" + string(suffix)
}

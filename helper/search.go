package helper

import "strings"

func NormalizeSearch(search string) string {
	return strings.TrimSpace(search)
}

func HasSearch(search string) bool {
	return NormalizeSearch(search) != ""
}

func SearchPattern(search string) string {
	return "%" + NormalizeSearch(search) + "%"
}

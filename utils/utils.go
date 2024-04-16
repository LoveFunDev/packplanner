package utils

import (
	"strconv"
	"strings"
)

// PrettyFormatFloat formats a float value to string
func PrettyFormatFloat(num float64, precise int) string {
	str := strconv.FormatFloat(num, 'f', precise, 64)
	if strings.Contains(str, ".") {
		return strings.TrimRight(strings.TrimRight(str, "0"), ".") // Trim right "0" and "."
	}

	return str
}

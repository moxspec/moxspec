package nvidia

import (
	"strconv"
	"strings"
)

func convUtilString(in string) float32 {
	return float32(parseFloatWithSuffix(in, "%"))
}

func convTempString(in string) float32 {
	return float32(parseFloatWithSuffix(in, "C"))
}

func convWattString(in string) float32 {
	return float32(parseFloatWithSuffix(in, "W"))
}

func convCountString(in string) int {
	in = strings.TrimSpace(in)
	if in == "" || in == "N/A" {
		return 0
	}
	val, _ := strconv.Atoi(in)
	return val
}

func parseFloatWithSuffix(in, suf string) float64 {
	in = strings.TrimSpace(in)
	if in == "" {
		return 0
	}
	if !strings.HasSuffix(in, suf) {
		log.Debugf("unexpexted input string (suffix '%s' not found in '%s')", suf, in)
		return 0
	}
	in = strings.TrimSuffix(in, suf)
	in = strings.TrimSpace(in)
	val, _ := strconv.ParseFloat(in, 64)
	return val
}

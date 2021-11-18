package util

import "math"

// Round returns rounded floting value
func Round(f float64, pt int) float64 {
	shift := math.Pow(10, float64(pt))
	return math.Floor(f*shift+.5) / shift
}

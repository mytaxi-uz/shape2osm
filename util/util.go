package util

import "math"

// truncate float64 value with precision 7 (~1cm)
func TruncateFloat64(num float64) float64 {
	const precesion = 10000000
	return float64(int(num*precesion+math.Copysign(0.5, num))) / precesion
}

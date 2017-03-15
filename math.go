package math0

import (
	"math"
)

const Epsilon float64 = 1.0e-8

func IsApproxEqual(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

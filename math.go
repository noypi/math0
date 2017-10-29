package math0

const Epsilon float64 = 1.0e-8
const Epsilon32 float32 = 1.0e-8

func IsApproxEqual(a, b float64) bool {
	if a > b {
		return (a - b) < Epsilon
	}

	return (b - a) < Epsilon
}

func IsApproxEqual32(a, b float32) bool {
	if a > b {
		return (a - b) < Epsilon32
	}

	return (b - a) < Epsilon32
}

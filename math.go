package math0

const Epsilon float64 = 1.0e-8

func IsApproxEqual(a, b float64) bool {
	if a > b {
		return (a - b) < Epsilon
	}

	return (b - a) < Epsilon
}

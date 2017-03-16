package math0

import (
	"math"
)

const Epsilon float64 = 1.0e-8

func IsApproxEqual(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

// TODO:support radical
func Power(a, power float64) float64 {
	if IsApproxEqual(0.0, power-(float64(int(power)))) {
		return powerBasic(a, power)
	}

	return powerRadical(a, power)
}

// https://en.wikipedia.org/wiki/Nth_root_algorithm
func powerRadical(a, power float64) float64 {
	return 0
}

// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
func powerBasic(a, power float64) float64 {
	if IsApproxEqual(math.Abs(power), 0.0) {
		return 1
	}

	if power < 0.0 {
		a = 1.0 / a
		power = -power
	}

	y := 1.0

	for 1 < power {
		if IsEven(power) {
			a = a * a
			power = power / 2.0
		} else {
			y = a * y
			a = a * a
			power = (power - 1) / 2
		}
	}
	return a * y

}

func IsEven(a float64) bool {
	return ((int(a) & 0x01) != 1) && (IsApproxEqual(0.0, a-(float64(int(a)))))
}

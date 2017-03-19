package kiwi

import (
	"fmt"
	"math"
)

func createStrength(a, b, c, w float64) StrengthType {
	result := 0.0
	result += math.Max(0.0, math.Min(1000.0, a*w)) * 1000000.0
	result += math.Max(0.0, math.Min(1000.0, b*w)) * 1000.0
	result += math.Max(0.0, math.Min(1000.0, c*w))
	return StrengthType(result)
}

func clipStrength(v StrengthType) StrengthType {
	return StrengthType(math.Max(0.0, math.Min(float64(_Required), float64(v))))
}

type StrengthType float64

var (
	_Required = createStrength(1000.0, 1000.0, 1000.0, 1.0)
	_Strong   = createStrength(1.0, 0.0, 0.0, 1.0)
	_Medium   = createStrength(0.0, 1.0, 0.0, 1.0)
	_Weak     = createStrength(0.0, 0.0, 1.0, 1.0)
)

func Required() StrengthType { return _Required }
func Strong() StrengthType   { return _Strong }
func Medium() StrengthType   { return _Medium }
func Weak() StrengthType     { return _Weak }

func (this StrengthType) String() string {
	switch this {
	case _Required:
		return "<strength:required>"
	case _Strong:
		return "<strength:strong>"
	case _Medium:
		return "<strength:medium>"
	case _Weak:
		return "<strength:weak>"
	}

	return fmt.Sprintf("<strength:%f>", float64(this))
}

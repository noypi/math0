package expr

import (
	"github.com/noypi/math0"
)

type Operator int
type TermConstructor func(c float64, varname string, power int) ITerm
type ExprConstructor func() IExpression
type EqnConstructor func() IEquation

const (
	OpEQ Operator = iota
	OpLEQ
	OpGEQ
	OpNEQ
	OpLess
	OpGreater
)

type IValuation interface {
	Get(varname string) (float64, bool)
}

func (this Operator) Test(a, b float64) bool {
	switch this {
	case OpEQ:
		return math0.IsApproxEqual(a, b)
	case OpLEQ:
		return a < b || math0.IsApproxEqual(a, b)
	case OpGEQ:
		return a > b || math0.IsApproxEqual(a, b)
	case OpNEQ:
		return !math0.IsApproxEqual(a, b)
	case OpLess:
		return a < b
	case OpGreater:
		return a > b
	}

	return false
}

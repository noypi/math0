package expr

import (
	"github.com/noypi/math0"
)

type Operator int
type TermConstructor func(c float64, varname string, power int) ITerm
type ExprConstructor func() IExpression
type EqnConstructor func() IEquation

const (
	OpLEQ Operator = iota
	OpGEQ
	OpEQ
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

func (this Operator) String() string {
	switch this {
	case OpLEQ:
		return "<="
	case OpGEQ:
		return ">="
	case OpEQ:
		return "=="
	case OpNEQ:
		return "!="
	case OpLess:
		return "<"
	case OpGreater:
		return ">"
	}
	return "<unkown operator>"
}

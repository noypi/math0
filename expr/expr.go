package expr

import (
	"github.com/noypi/math0"
)

type Relation int
type TermConstructor func(c float64, varname string, power int) ITerm
type ExprConstructor func() IExpression
type EqnConstructor func() IEquation

const (
	LEQ Relation = iota
	GEQ
	EQ
	NEQ
	Lesser
	Greater
)

type IValuation interface {
	Get(varname string) (float64, bool)
}

func (this Relation) Test(a, b float64) bool {
	switch this {
	case EQ:
		return math0.IsApproxEqual(a, b)
	case LEQ:
		return a < b || math0.IsApproxEqual(a, b)
	case GEQ:
		return a > b || math0.IsApproxEqual(a, b)
	case NEQ:
		return !math0.IsApproxEqual(a, b)
	case Lesser:
		return a < b
	case Greater:
		return a > b
	}

	return false
}

func (this Relation) String() string {
	switch this {
	case LEQ:
		return "<="
	case GEQ:
		return ">="
	case EQ:
		return "=="
	case NEQ:
		return "!="
	case Lesser:
		return "<"
	case Greater:
		return ">"
	}
	return "<unkown relation>"
}

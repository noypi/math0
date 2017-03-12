package expr

import (
	"sort"

	. "github.com/noypi/math0"
)

type Operator int
type TermConstructor func(c float64, varname string, power int) ITerm
type ExprConstructor func() IExpression
type EqnConstructor func() IEquation

var (
	NewTerm TermConstructor
	NewExpr ExprConstructor
	NewEqn  EqnConstructor
)

const (
	OpEQ Operator = iota
	OpLEQ
	OpGEQ
	OpNEQ
	OpLess
	OpGreater
)

type IEquation interface {
	Left() IExpression
	Right() IExpression
	Op() Operator

	SetLeft(IExpression)
	SetRight(IExpression)
	SetOp(Operator)
}

type IExpression interface {
	EachTerm(func(ITerm) bool)
	EachTermByVar(func(varname string, power int, xs TermList) bool)
	Terms() TermList   // sorted by power
	AddTerm(...ITerm)  // appends
	SetTerms(...ITerm) // clears, then sets
}

type TermList []ITerm

type ITerm interface {
	C() float64
	Var() string
	Power() int
	SetC(c float64)
	SetVar(x string)
	SetPower(n int)
}

type IValuation interface {
	Get(varname string, defaultValue float64) float64
}

func (this TermList) SortByPower() {
	sort.Slice(this, func(i, j int) bool {
		return this[i].Power() >= this[j].Power()
	})
}

func CopyTerm(term ITerm) ITerm {
	return NewTerm(term.C(), term.Var(), term.Power())
}

func CopyExpr(expr IExpression) IExpression {
	o := NewExpr()
	expr.EachTerm(func(term ITerm) bool {
		o.AddTerm(CopyTerm(term))
		return true
	})
	return o
}

func CopyEqn(eqn IEquation) IEquation {
	o := NewEqn()
	o.SetLeft(CopyExpr(eqn.Left()))
	o.SetRight(CopyExpr(eqn.Right()))
	o.SetOp(eqn.Op())
	return o
}

func SimplifyExpression(expr IExpression) {
	var terms []ITerm // new terms
	var bModified bool

	expr.EachTermByVar(func(varname string, power int, xs TermList) bool {
		if 2 <= len(xs) {
			bModified = true
		} else if 1 == len(xs) {
			terms = append(terms, xs[0])
		} else if 0 == len(xs) {
			return true
		}

		var total float64
		for _, term := range xs {
			total += term.C()
		}

		if IsApproxEqual(total, 0.0) {
			return true
		}

		terms = append(terms, NewTerm(total, varname, power))
		return true
	})

	if bModified {
		expr.SetTerms(terms...)
	}
}

func ValueOfExpr(expr IExpression, m IValuation) float64 {
	var total float64
	expr.EachTerm(func(term ITerm) bool {
		total += term.C() * (m.Get(term.Var(), 0.0))
		return true
	})

	return total
}

func IsEquationTrue(eqn IEquation, m IValuation) bool {
	left := ValueOfExpr(eqn.Left(), m)
	right := ValueOfExpr(eqn.Right(), m)
	return eqn.Op().Test(left, right)
}

func (this Operator) Test(a, b float64) bool {
	switch this {
	case OpEQ:
		return IsApproxEqual(a, b)
	case OpLEQ:
		return a < b || IsApproxEqual(a, b)
	case OpGEQ:
		return a > b || IsApproxEqual(a, b)
	case OpNEQ:
		return !IsApproxEqual(a, b)
	case OpLess:
		return a < b
	case OpGreater:
		return a > b
	}

	return false
}

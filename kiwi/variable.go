package kiwi

import (
	"github.com/noypi/math0/expr"
)

type Variable struct {
	expr.IVariable
	value float64
}

func Var(name string) *Variable {
	o := new(Variable)
	o.IVariable = expr.NewVar(name)
	return o
}

func ValueOf(v expr.IVariable) float64 {
	return v.(*Variable).value
}

func (this Variable) Value() float64 {
	return this.value
}

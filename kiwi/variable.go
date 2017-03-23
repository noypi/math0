package kiwi

import (
	"github.com/noypi/math0/expr"
)

type _Variable struct {
	expr.IVariable
	value float64
}

func Variable(name string) *_Variable {
	o := new(_Variable)
	o.IVariable = expr.NewVar(name)
	return o
}

func ValueOf(v expr.IVariable) float64 {
	return v.(*_Variable).value
}

func (this _Variable) Value() float64 {
	return this.value
}

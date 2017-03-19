package kiwi

import (
	"bytes"
	"fmt"

	"github.com/noypi/math0/expr"
)

/*
type RelationalOperator int

const (
	OP_LE RelationalOperator = iota
	OP_GE
	OP_EQ
)*/

type _Constraint struct {
	strength   StrengthType
	expression expr.IExpression
	op         expr.Operator
}

func Constraint(eqn expr.IEquation, strength StrengthType) *_Constraint {
	o := new(_Constraint)
	o.strength = strength
	o.op = eqn.Op()

	var terms expr.TermList
	eqn.Right().EachTerm(func(term expr.ITerm) bool {
		terms = append(terms, expr.NewTerm(-term.C(), term.Vars()...))
		return true
	})
	eqn.Left().AddTerm(terms...)
	eqn.Right().AddTerm(terms...)
	o.expression = eqn.Left()
	return o
}

func (this _Constraint) Constant() float64 {
	return this.expression.Constant()
}

func (this _Constraint) String() string {
	return fmt.Sprintf("%s,%s", this.strength, this.expression)
}

func (this _Constraint) Dump(buf *bytes.Buffer) {
	this.expression.EachTerm(func(term expr.ITerm) bool {
		if 0 == len(term.Vars()) {
			return true
		}

		buf.WriteString(fmt.Sprintf("%s * %v + ", expr.ToTrimZero(term.C()), term.VarAt(0).Name()))
		return true
	})
	buf.WriteString(expr.ToTrimZero(this.expression.Constant()))

	switch this.op {
	case expr.OpLEQ:
		buf.WriteString(" <= 0 ")
	case expr.OpGEQ:
		buf.WriteString(" >= 0 ")
	case expr.OpEQ:
		buf.WriteString(" == 0 ")
	}

	buf.WriteString(" | strength = ")
	buf.WriteString(this.strength.String())
	buf.WriteString("\n")
}

/*
func (this RelationalOperator) String() string {
	switch this {
	case OP_LE:
		return "OP_LE"
	case OP_GE:
		return "OP_GE"
	case OP_EQ:
		return "OP_EQ"
	}

	return "<unknown op>"
}
*/

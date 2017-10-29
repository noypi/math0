package kiwi

import (
	"bytes"
	"fmt"

	"github.com/noypi/math0/expr"
)

type Constraint struct {
	strength   StrengthType
	expression expr.IExpression
	relation   expr.Relation
}

func NewConstraint(eqn expr.IEquation, strength StrengthType) *Constraint {
	o := new(Constraint)
	o.strength = strength
	o.relation = eqn.Relation()
	expr2 := eqn.Left().Clone()

	eqn.Right().EachTerm(func(term expr.ITerm) bool {
		expr2.AddTerm(expr.NewTerm(-term.C(), term.Vars()...))
		return true
	})
	o.expression = expr2
	return o
}

func (this Constraint) Constant() float64 {
	return this.expression.Constant()
}

func (this Constraint) Expression() expr.IExpression {
	return this.expression
}

func (this Constraint) String() string {
	return fmt.Sprintf("%s,%s", this.strength, this.expression)
}

func (this Constraint) Dump() string {
	buf := bytes.NewBufferString("")
	this.expression.EachTerm(func(term expr.ITerm) bool {
		if 0 == len(term.Vars()) {
			return true
		}

		buf.WriteString(fmt.Sprintf("%s*%v + ", expr.ToTrimZero(term.C()), term.VarAt(0).Name()))
		return true
	})
	buf.WriteString(expr.ToTrimZero(this.expression.Constant()))

	switch this.relation {
	case expr.LEQ:
		buf.WriteString(" <= 0 ")
	case expr.GEQ:
		buf.WriteString(" >= 0 ")
	case expr.EQ:
		buf.WriteString(" == 0 ")
	}

	buf.WriteString(" | strength = ")
	buf.WriteString(this.strength.String())
	buf.WriteString("\n")

	return buf.String()
}

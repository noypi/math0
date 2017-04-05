package expr

import (
	"fmt"
)

type IEquation interface {
	fmt.Stringer
	Left() IExpression
	Right() IExpression
	Relation() Relation

	SetLeft(IExpression)
	SetRight(IExpression)
	SetRelation(Relation)
	Clone() IEquation
}

func Equation(left IExpression, rel Relation, right IExpression) IEquation {
	o := new(_Equation)
	o.left = left
	o.relation = rel
	o.right = right

	if nil == o.left {
		o.left = NewExpr()
		o.left.AddTerm(NewTerm(0))
	}
	if nil == o.right {
		o.right = NewExpr()
		o.right.AddTerm(NewTerm(0))
	}

	return o
}

func IsEquationTrue(eqn IEquation, m IValuation) (b bool, err error) {
	var left, right float64
	if left, err = ValueOfExpr(eqn.Left(), m); nil != err {
		return
	}
	if right, err = ValueOfExpr(eqn.Right(), m); nil != err {
		return
	}
	b = eqn.Relation().Test(left, right)
	return
}

type _Equation struct {
	left     IExpression
	right    IExpression
	relation Relation
}

func (this _Equation) Left() IExpression {
	return this.left
}

func (this _Equation) Right() IExpression {
	return this.right
}

func (this *_Equation) Relation() Relation {
	return this.relation
}

func (this *_Equation) SetLeft(expr IExpression) {
	this.left = expr
}

func (this *_Equation) SetRight(expr IExpression) {
	this.right = expr
}

func (this *_Equation) SetRelation(relation Relation) {
	this.relation = relation
}

func (this _Equation) Clone() IEquation {
	o := new(_Equation)
	o.SetLeft(this.left.Clone())
	o.SetRight(this.right.Clone())
	o.SetRelation(this.relation)
	return o
}

func (this _Equation) String() string {
	return fmt.Sprintf("%s %s %s", this.left, this.relation, this.right)
}

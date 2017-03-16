package expr

import (
	"fmt"
	"math"

	"github.com/noypi/math0"
)

type IExpression interface {
	EachTerm(func(ITerm) bool)
	EachTermByVar(vars VariableList, cb func(ITerm) bool)
	AllTermsWithVar(func(vars VariableList, xs TermList) bool)
	Terms() TermList   // sorted by power
	AddTerm(...ITerm)  // appends
	SetTerms(...ITerm) // clears, then sets
	Clone() IExpression
}

func ValueOfExpr(expr IExpression, m IValuation) (value float64, err error) {
	expr.EachTerm(func(term ITerm) bool {
		for _, v := range term.Vars() {
			c, has := m.Get(v.Name())
			if !has {
				err = ErrNoValuationForVar(fmt.Errorf("var=%s has no valuation", v.Name()))
				return false
			}
			if math0.IsApproxEqual(v.Power(), 0.0) {
				value += term.C()
			} else {
				value += term.C() * math.Pow(c, v.Power())
			}

		}
		return true
	})

	return
}

func SimplifyExpression(expr IExpression) {
	var terms []ITerm // new terms
	var bModified bool

	expr.AllTermsWithVar(func(vars VariableList, xs TermList) bool {
		if 2 <= len(xs) {
			bModified = true
		} else if 0 == len(xs) {
			panic("should not be here")
		}

		var total float64
		for _, term := range xs {
			total += term.C()
		}

		if math0.IsApproxEqual(total, 0.0) {
			bModified = true
			return true
		}

		o := xs[0].Clone()
		o.SetC(total)
		terms = append(terms, o)
		return true
	})

	if bModified {
		expr.SetTerms(terms...)
	}
}

type _termsMapVal struct {
	termlist TermList
	vars     VariableList
}
type _Expression struct {
	terms map[string]_termsMapVal
}

func NewExpr() IExpression {
	return &_Expression{
		terms: map[string]_termsMapVal{},
	}
}

func (this *_Expression) AddTerm(terms ...ITerm) {
	for _, term := range terms {
		k := term.Vars().String()
		mapVal, has := this.terms[k]
		if !has {
			mapVal.vars = term.Vars()
		}
		mapVal.termlist = append(mapVal.termlist, term)
		this.terms[k] = mapVal
	}
}

func (this *_Expression) SetTerms(terms ...ITerm) {
	this.terms = map[string]_termsMapVal{}
	this.AddTerm(terms...)
}

func (this _Expression) AllTermsWithVar(cb func(vars VariableList, xs TermList) bool) {
	for _, vals := range this.terms {
		if 0 == len(vals.termlist) {
			continue
		}
		if !cb(vals.vars, vals.termlist) {
			return
		}
	}
}

func (this _Expression) EachTermByVar(vars VariableList, cb func(term ITerm) bool) {
	vals, _ := this.terms[vars.String()]
	for _, term := range vals.termlist {
		if !cb(term) {
			return
		}
	}
}

func (this _Expression) Terms() TermList {
	var ls TermList
	for _, vals := range this.terms {
		ls = append(ls, vals.termlist...)
	}
	return ls
}

func (this _Expression) EachTerm(cb func(term ITerm) bool) {
	for _, vals := range this.terms {
		for _, term := range vals.termlist {
			if !cb(term) {
				return
			}
		}
	}
}

func (this _Expression) Clone() IExpression {
	o := NewExpr()
	this.EachTerm(func(term ITerm) bool {
		o.AddTerm(term.Clone())
		return true
	})
	return o
}

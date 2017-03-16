package expr

import (
	"bytes"
	"fmt"
	"math"
	"sort"

	"github.com/noypi/math0"
)

type IExpression interface {
	fmt.Stringer
	EachTerm(func(ITerm) bool)
	Terms() TermList   // sorted by power
	AddTerm(...ITerm)  // appends
	SetTerms(...ITerm) // clears, then sets
	Key() string
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

func SimplifyExpression(expr IExpression) (out TermList, bDidSomething bool) {

	terms := expr.Terms()
	if 1 >= len(terms) {
		return
	}
	if terms.IsSimplified() {
		return
	}
	if !terms.IsSorted() {
		sort.Slice(terms, terms.Less)
	} else {
		bDidSomething = true
	}

	out = TermList{terms[0]}
	outPrev := terms[0]
	for i := 1; i < len(terms); i++ {
		if math0.IsApproxEqual(terms[i].C(), 0.0) {
			continue
		}

		if terms[i].Key() == outPrev.Key() {
			outPrev.SetC(outPrev.C() + terms[i].C())
			if math0.IsApproxEqual(outPrev.C(), 0.0) {
				out = out[:len(out)-1]
			}
		} else {
			out = append(out, terms[i])
			outPrev = terms[i]
		}
	}

	return
}

type _Expression struct {
	terms TermList
	key   *string
}

func NewExpr() IExpression {
	return &_Expression{
		terms: TermList{},
	}
}

func (this *_Expression) Key() string {
	if nil != this.key {
		return *this.key
	}
	if 0 == len(this.terms) {
		k := "0"
		this.key = &k
		return *this.key
	}

	buf := bytes.NewBufferString(this.terms[0].Key())
	if 1 == len(this.terms) {
		return buf.String()
	}

	for _, term := range this.terms[1:] {
		buf.WriteString(",")
		buf.WriteString(term.Key())
	}

	k := buf.String()
	this.key = &k
	return *this.key
}

func (this _Expression) String() string {
	if 0 == len(this.terms) {
		return "0"
	}

	buf := bytes.NewBufferString(this.terms[0].String())
	if 1 == len(this.terms) {
		return buf.String()
	}

	for _, term := range this.terms[1:] {
		if 0 < term.C() {
			buf.WriteString(" + ")
		} else {
			buf.WriteString(" - ")
		}
		buf.WriteString(term.String())
	}

	return buf.String()
}

func (this *_Expression) AddTerm(terms ...ITerm) {
	this.key = nil
	this.terms = append(this.terms, terms...)
	if newTerms, bWasSimplified := SimplifyExpression(this); bWasSimplified {
		this.terms = newTerms
	}
}

func (this *_Expression) SetTerms(terms ...ITerm) {
	// this.key = nil // is set in AddTerm()
	this.terms = this.terms[0:]
	this.AddTerm(terms...)
}

func (this TermList) IsSorted() bool {
	if 1 >= len(this) {
		return true
	}
	return sort.SliceIsSorted(this, this.Less)
}

func (this TermList) IsSimplified() bool {
	if 1 >= len(this) {
		return true
	}
	for iprev, i := 0, 1; i < len(this); iprev, i = iprev+1, i+1 {
		if !this.Less(iprev, i) {
			return false
		}
	}
	return true
}

func (this TermList) Less(i, j int) bool {
	if this[i].PowerTotal() == this[j].PowerTotal() {
		return this[i].Key() < this[j].Key()
	}
	return this[i].PowerTotal() < this[j].PowerTotal()
}

func (this _Expression) Terms() TermList {
	return this.terms
}

func (this _Expression) EachTerm(cb func(term ITerm) bool) {
	for _, term := range this.terms {
		if !cb(term) {
			return
		}
	}
}

func (this _Expression) Clone() IExpression {
	o := NewExpr()
	o.SetTerms(this.terms...)
	return o
}

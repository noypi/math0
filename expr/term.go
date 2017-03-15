package expr

import (
	"sort"
)

type ITerm interface {
	C() float64
	SetC(c float64)

	Var(name string) (v IVariable)
	Vars() VariableList
	SetVars(vs ...IVariable)
	VarAt(i int) IVariable

	Clone() ITerm
}

type TermList []ITerm

type _Term struct {
	c    float64
	vars VariableList
}

func NewTerm(c float64, vs ...IVariable) ITerm {
	o := &_Term{c: c}
	o.SetVars(vs...)
	return o
}

func (this _Term) C() float64 {
	return this.c
}

func (this _Term) VarAt(i int) (v IVariable) {
	if i < len(this.vars) && 0 <= i {
		return this.vars[i]
	}
	return nil
}

func (this _Term) Vars() VariableList {
	return this.vars
}

func (this _Term) Var(name string) (v IVariable) {
	i := sort.Search(len(this.vars), func(i int) bool {
		return this.vars[i].Name() >= name
	})
	if i < len(this.vars) {
		v = this.vars[i]
	}

	return nil
}

func (this _Term) SetC(c float64) {
	this.c = c
}

func (this _Term) SetVars(vs ...IVariable) {
	this.vars = this.vars[:0]
	if 0 == len(vs) {
		return
	}

	this.vars = append(this.vars, vs[0])
	if 1 == len(vs) {
		return
	}

	for i, _ := range vs[1:] {
		jprev := len(this.vars) - 1
		if vs[i].Name() == this.vars[jprev].Name() {
			this.vars[jprev].AddPower(vs[i].Power())
		} else {
			this.vars = append(this.vars, vs[i])
		}
	}
	if 1 < len(vs) {
		sort.Slice(vs, func(i, j int) bool {
			return vs[i].Name() < vs[j].Name()
		})
	}
	this.vars = vs
}

func (this _Term) Clone() ITerm {
	return NewTerm(this.c, this.vars...)
}

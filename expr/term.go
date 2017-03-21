package expr

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/noypi/math0"
)

type ITerm interface {
	fmt.Stringer
	C() float64
	SetC(c float64)

	Var(name string) (v IVariable)
	Vars() VariableList
	SetVars(vs ...IVariable)
	VarAt(i int) IVariable

	Clone() ITerm
	Key() string
	PowerTotal() float64
}

type TermList []ITerm

type _Term struct {
	c          float64
	vars       VariableList
	key        *string
	powertotal *float64
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
	if i < len(this.vars) && this.vars[i].Name() == name {
		v = this.vars[i]
	}

	return
}

func (this *_Term) SetC(c float64) {
	this.c = c
}

func (this *_Term) SetVars(vs ...IVariable) {
	this.key = nil
	this.powertotal = nil
	this.vars = this.vars[:0]
	if 0 == len(vs) {
		return
	}

	if !math0.IsApproxEqual(vs[0].Power(), 0.0) {
		this.vars = append(this.vars, vs[0])
	}
	if 1 == len(vs) {
		return
	}

	for i, _ := range vs[1:] {
		if math0.IsApproxEqual(vs[i].Power(), 0.0) {
			continue
		}

		jprev := len(this.vars) - 1
		if vs[i].Name() == this.vars[jprev].Name() {
			this.vars[jprev] = this.vars[jprev].AddPower(vs[i].Power())
			if math0.IsApproxEqual(this.vars[jprev].Power(), 0.0) {
				// removed if 0.0
				this.vars = this.vars[:jprev]
			}
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

func (this *_Term) Key() string {
	if nil == this.key {
		k := this.vars.Key()
		this.key = &k
	}
	return *this.key
}

func (this _Term) String() string {
	if 0 == len(this.vars) {
		return ToTrimZero(this.c)
	}
	return fmt.Sprintf("%s(%s)", ToTrimZero(this.c), this.vars.String())
}

func (this *_Term) PowerTotal() float64 {
	if nil == this.powertotal {
		n := this.vars.PowerTotal()
		this.powertotal = &n
	}
	return *this.powertotal
}

func (this _Term) Clone() ITerm {
	return NewTerm(this.c, this.vars...)
}

func (this TermList) String() string {
	if 0 == len(this) {
		return "0"
	}

	buf := bytes.NewBufferString(this[0].String())
	if 1 == len(this) {
		return buf.String()
	}

	for _, term := range this[1:] {
		buf.WriteString(" + ")
		buf.WriteString(term.String())
	}

	return buf.String()
}

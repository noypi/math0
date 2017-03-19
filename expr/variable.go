package expr

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/noypi/math0"
)

// should be immutable
type IVariable interface {
	fmt.Stringer
	Name() string
	Power() float64
	AddPower(float64) IVariable
}

type _Variable struct {
	name  string
	power float64
}

type VariableList []IVariable

func NewVarN(name string, power float64) IVariable {
	if 1 <= power {
		return &_Variable{name: name, power: power}
	}
	return &_Variable{}
}

func NewVar(name string) IVariable {
	return NewVarN(name, 1.0)
}

func (this _Variable) Name() string {
	return this.name
}

func (this _Variable) Power() float64 {
	return this.power
}

func (this *_Variable) AddPower(n float64) IVariable {
	return NewVarN(this.name, this.power+n)
}

func (this _Variable) String() string {
	if math0.IsApproxEqual(this.power, 1.0) {
		return this.name
	}

	return fmt.Sprintf("%s^%s", this.name, ToTrimZero(this.power))
}

func (this VariableList) Simplify() (out VariableList) {
	if this.IsSimplified() {
		return this
	}

	vs := this
	if 1 >= len(vs) {
		return this
	}

	if !vs.IsSorted() {
		sort.Slice(vs, vs.Less)
	}

	out = VariableList{vs[0]}
	outPrev := vs[0]
	for i := 1; i < len(vs); i++ {
		if math0.IsApproxEqual(vs[i].Power(), 0.0) {
			continue
		}

		if outPrev.Name() == vs[i].Name() {
			outPrev = outPrev.AddPower(vs[i].Power())
			if math0.IsApproxEqual(outPrev.Power(), 0.0) {
				out = out[:len(out)-1]
			} else {
				out[len(out)-1] = outPrev
			}
		} else {
			if !math0.IsApproxEqual(vs[i].Power(), 0.0) {
				out = append(out, vs[i])
				outPrev = out[len(out)-1]
			}
		}

	}

	return out
}

func (this VariableList) Less(i, j int) bool {
	return this[i].Name() < this[j].Name()
}

func (this VariableList) IsSimplified() bool {
	if 1 >= len(this) {
		return true
	}

	for i, iprev := 1, 0; i < len(this); i, iprev = i+1, iprev+1 {
		if !this.Less(iprev, i) {
			return false
		}
	}
	return true
}

func (this VariableList) IsSorted() bool {
	if 1 >= len(this) {
		return true
	}
	return sort.SliceIsSorted(this, this.Less)
}

func (this *VariableList) Key() string {
	if 1 < len(*this) {
		*this = this.Simplify()
	}
	return (*this).String()
}

func (this VariableList) PowerTotal() float64 {
	if 0 == len(this) {
		return 0.0
	} else if 1 == len(this) {
		return this[0].Power()
	}

	total := 0.0
	for _, v := range this {
		total += v.Power()
	}
	return total
}

func (this VariableList) String() string {
	if nil == this || 0 == len(this) {
		return ""
	} else if 1 == len(this) {
		return this[0].String()
	}

	buf := bytes.NewBufferString("")
	for _, v := range this {
		if 0 < buf.Len() {
			buf.WriteString("*")
		}
		if math0.IsApproxEqual(v.Power(), 1.0) {
			buf.WriteString(v.Name())
		} else {
			buf.WriteString(v.String())
		}
	}

	return buf.String()
}

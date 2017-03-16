package expr

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/noypi/math0"
)

type IVariable interface {
	fmt.Stringer
	Name() string
	Power() float64
	AddPower(float64) float64
}

type _Variable struct {
	name  string
	power float64
}

type VariableList []IVariable

func VariableN(name string, power float64) IVariable {
	if 1 <= power {
		return &_Variable{name: name, power: power}
	}
	return &_Variable{}
}

func Variable(name string) IVariable {
	return VariableN(name, 1.0)
}

func (this _Variable) Name() string {
	return this.name
}

func (this _Variable) Power() float64 {
	return this.power
}

func (this *_Variable) AddPower(n float64) float64 {
	this.power += n
	return this.power
}

func (this _Variable) String() string {
	spow := fmt.Sprintf("%f", this.power)
	spow = strings.TrimRight(spow, ".0")

	if math0.IsApproxEqual(this.power, 1.0) {
		return this.name
	}

	return fmt.Sprintf("%s^%s", this.name, spow)
}

func (this VariableList) Simplify() (out VariableList) {
	vs := this
	if 1 >= len(vs) {
		return this
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Name() < vs[j].Name()
	})

	out = VariableList{vs[0]}
	outPrev := vs[0]
	for i := 1; i < len(vs); i++ {
		if math0.IsApproxEqual(vs[i].Power(), 0.0) {
			continue
		}

		if outPrev.Name() == vs[i].Name() {
			outPrev.AddPower(vs[i].Power())
			if math0.IsApproxEqual(outPrev.Power(), 0.0) {
				out = out[:len(out)-1]
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

func (this *VariableList) Key() string {
	*this = this.Simplify()
	return (*this).String()
}

func (this VariableList) String() string {

	if 0 == len(this) {
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

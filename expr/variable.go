package expr

import (
	"bytes"
	"fmt"

	"github.com/noypi/mapk"
	"github.com/noypi/math0"
)

type IVariable interface {
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

func (this VariableList) String() string {
	if 0 == len(this) {
		return ""
	} else if 1 == len(this) {
		return this[0].Name()
	}

	m := mapk.Map(mapk.CmpString)
	for _, v := range this {
		n := m.Get(v.Name())
		if nil == n {
			m.Put(v.Name(), 1.0)
		} else {
			m.Put(v.Name(), n.(float64)+v.Power())
		}
	}

	buf := bytes.NewBufferString("")
	m.EachFrom("", func(name, power interface{}) bool {
		if 0 < buf.Len() {
			buf.WriteString("*")
		}
		if math0.IsApproxEqual(power.(float64), 1.0) {
			buf.WriteString(name.(string))
		} else {
			// TODO // trim zeros
			buf.WriteString(fmt.Sprintf("%s^%d", name, power))
		}
		return true
	})

	return buf.String()
}

package expr

import (
	"bytes"
	"fmt"

	"github.com/noypi/mapk"
)

type IVariable interface {
	Name() string
	Power() int
	AddPower(int) int
}

type _Variable struct {
	name  string
	power int
}

type VariableList []IVariable

func VariableN(name string, power int) IVariable {
	if 1 <= power {
		return &_Variable{name: name, power: power}
	}
	return &_Variable{}
}

func Variable(name string) IVariable {
	return VariableN(name, 1)
}

func (this _Variable) Name() string {
	return this.name
}

func (this _Variable) Power() int {
	return this.power
}

func (this *_Variable) AddPower(n int) int {
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
			m.Put(v.Name(), 1)
		} else {
			m.Put(v.Name(), n.(int)+1)
		}
	}

	buf := bytes.NewBufferString("")
	m.EachFrom("", func(name, power interface{}) bool {
		if 0 < buf.Len() {
			buf.WriteString("*")
		}
		if 1 < power.(int) {
			buf.WriteString(fmt.Sprintf("%s^%d", name, power))
		} else {
			buf.WriteString(name.(string))
		}
		return true
	})

	return buf.String()
}

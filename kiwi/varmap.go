package kiwi

import (
	"bytes"

	"github.com/noypi/math0/expr"
)

type _varMapVal struct {
	k *_Variable
	v _Symbol
}
type _VarMap map[string]*_varMapVal

func (this _VarMap) Get(v expr.IVariable) (symbol _Symbol, has bool) {
	if nil == v {
		return
	}
	o, has := this[v.Name()]
	if has {
		symbol = o.v
	}
	return
}

func (this _VarMap) Put(v expr.IVariable, symbol _Symbol) {
	this[v.Name()] = &_varMapVal{k: v.(*_Variable), v: symbol}
}

func (this _VarMap) Each(cb func(expr.IVariable, _Symbol) bool) {
	for _, v := range this {
		if !cb(v.k, v.v) {
			break
		}
	}
}

func (this _VarMap) Dump(buf *bytes.Buffer) {
	for _, v := range this {
		buf.WriteString(v.k.Name() + " = ")
		v.v.Dump(buf)
		buf.WriteString("\n")
	}
}

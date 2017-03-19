package kiwi

import (
	"github.com/noypi/math0/expr"
)

type _EditMap map[string]*_EditInfo

func (this _EditMap) Get(v expr.IVariable) (info *_EditInfo, has bool) {
	info, has = this[v.Name()]
	return
}

func (this _EditMap) Put(v expr.IVariable, info *_EditInfo) {
	this[v.Name()] = info
}

func (this _EditMap) Delete(v expr.IVariable) {
	delete(this, v.Name())
}

package kiwi

import (
	"bytes"
)

type _CnMap map[string]*_cnMapVal
type _cnMapVal struct {
	k *Constraint
	v *_Tag
}

func (this _CnMap) Get(cn *Constraint) (tag *_Tag, has bool) {
	v, has := this[cn.String()]
	if has {
		tag = v.v
	}
	return
}

func (this _CnMap) Put(cn *Constraint, tag *_Tag) {
	this[cn.String()] = &_cnMapVal{k: cn, v: tag}
}

func (this _CnMap) Delete(cn *Constraint) {
	delete(this, cn.String())
}

func (this _CnMap) Dump() string {
	buf := bytes.NewBufferString("")
	for _, v := range this {
		buf.WriteString(v.k.Dump())
	}
	return buf.String()
}

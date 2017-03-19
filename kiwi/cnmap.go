package kiwi

import (
	"bytes"
)

type _CnMap map[string]*_cnMapVal
type _cnMapVal struct {
	k *_Constraint
	v *_Tag
}

func (this _CnMap) Get(cn *_Constraint) (tag *_Tag, has bool) {
	v, has := this[cn.String()]
	if has {
		tag = v.v
	}
	return
}

func (this _CnMap) Put(cn *_Constraint, tag *_Tag) {
	this[cn.String()] = &_cnMapVal{k: cn, v: tag}
}

func (this _CnMap) Delete(cn *_Constraint) {
	delete(this, cn.String())
}

func (this _CnMap) Dump(buf *bytes.Buffer) {
	for _, v := range this {
		v.k.Dump(buf)
	}
}

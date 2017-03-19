package kiwi

import (
	"bytes"
)

type _RowMap map[_Symbol]*_Row

func (this _RowMap) Get(sym _Symbol) (row *_Row, has bool) {
	row, has = this[sym]
	return
}

func (this _RowMap) Put(sym _Symbol, row *_Row) {
	this[sym] = row
}

func (this _RowMap) Delete(sym _Symbol) {
	delete(this, sym)
}

func (this _RowMap) Dump() string {
	buf := bytes.NewBufferString("")
	for k, v := range this {
		buf.WriteString(k.Dump())
		buf.WriteString(" | ")
		buf.WriteString(v.Dump())
	}

	return buf.String()
}

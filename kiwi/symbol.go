package kiwi

import (
	"bytes"
	"fmt"
)

type _Symbol struct {
	Id   int64
	Type SymbolType
}

type SymbolType int

const (
	Invalid SymbolType = iota
	External
	Slack
	Error
	Dummy
)

func (this _Symbol) Dump(buf *bytes.Buffer) {
	switch this.Type {
	case Invalid:
		buf.WriteString("i")
	case External:
		buf.WriteString("v")
	case Slack:
		buf.WriteString("s")
	case Error:
		buf.WriteString("e")
	case Dummy:
		buf.WriteString("d")
	}
	buf.WriteString(fmt.Sprintf("%d", this.Id))
}

type _SymbolList []_Symbol

func (this _SymbolList) Dump(buf *bytes.Buffer) {
	for _, sym := range this {
		sym.Dump(buf)
		buf.WriteString("\n")
	}
}

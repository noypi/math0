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

func (this _Symbol) Dump() string {
	buf := bytes.NewBufferString("")
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

	return buf.String()
}

type _SymbolList []_Symbol

func (this _SymbolList) Dump() string {
	buf := bytes.NewBufferString("")
	for _, sym := range this {
		buf.WriteString(sym.Dump())
		buf.WriteString("\n")
	}
	return buf.String()
}

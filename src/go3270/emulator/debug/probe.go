package debug

import (
	"fmt"
	"go3270/emulator/conv"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logProbe(addr int) {
	t := l.newTable(text.FgHiMagenta, "")
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "Row", "Col", "SF", "Blink", "Color", "Hidden", "Hilite", "MDT", "Num", "Prot", "Rev", "Uscore"})
	cell, _ := l.buf.Peek(addr)
	l.withAttrs(t, fmt.Sprintf("%#02x '%s'", conv.A2E(cell.Char), utils.Ternary(cell.Char >= 0x20, string(cell.Char), " ")), addr, cell.Attrs, cell.FldStart)
}

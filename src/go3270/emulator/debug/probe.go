package debug

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logProbe(addr int) {
	t := l.newTable(text.FgHiMagenta)
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "Row", "Col", "SF", "Blink", "Color", "Hidden", "Hilite", "MDT", "Num", "Prot", "Rev", "Uscore"})
	cell, _ := l.buf.Peek(addr)
	l.withAttrs(t, "Probe", addr, cell.Attrs, cell.FldStart)
}

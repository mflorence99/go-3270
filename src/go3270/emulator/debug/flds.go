package debug

import (
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log buffer fields

func (l *Logger) logFlds() {
	t := l.newTable(text.FgHiBlue, "Buffer Fields")
	defer t.Render()
	// 👇 header rows
	t.AppendHeader(table.Row{"Row", "Col", "Len", "PROT", "MDT", "Data"})
	// 👇 data rows
	for _, fld := range l.flds.Get() {
		sf, _ := fld.StartFld()
		row, col := l.cfg.Addr2RC(sf.FldAddr)
		// 👇 gather all the chars in the fld
		t.AppendRow(table.Row{
			row,
			col,
			len(fld),
			utils.Ternary(sf.Attrs.Protected, "PROT", ""),
			utils.Ternary(sf.Attrs.Modified, "MDT", ""),
			utils.Truncate(fld.String(), 60),
		})
	}
}

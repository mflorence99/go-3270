package debug

import (
	"go3270/emulator/buffer"
	"go3270/emulator/utils"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logFlds(flds *buffer.Flds) {
	t := l.newTable(text.FgHiBlue, "Flds (before render)")
	defer t.Render()
	// 👇 header rows
	t.AppendHeader(table.Row{"Row", "Col", "PROT", "MDT", "Data"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 5, Transformer: l.wrap(80), WidthMax: 80},
	})
	// 👇 data rows
	for _, fld := range flds.Get() {
		sf, _ := fld.StartFld()
		row, col := l.cfg.Addr2RC(sf.FldAddr)
		data := make([]byte, len(fld)-1)
		for ix := 1; ix < len(fld); ix++ {
			cell := fld[ix]
			data[ix-1] = utils.Ternary(cell.Char <= ' ', ' ', cell.Char)
		}
		t.AppendRow(table.Row{
			row,
			col,
			utils.Ternary(sf.Attrs.Protected, "PROT", ""),
			utils.Ternary(sf.Attrs.Modified, "MDT", ""),
			strings.TrimSpace(string(data)),
		})
	}
}

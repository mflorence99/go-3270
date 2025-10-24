package debug

import (
	"go3270/emulator/wcc"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logWCC(wcc wcc.WCC) {
	t := l.newTable()
	defer t.Render()
	// 👇 table rows
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter},
		{Number: 4, Align: text.AlignCenter},
		{Number: 5, Align: text.AlignCenter},
	})
	t.AppendHeader(table.Row{"", "Alarm", "Reset", "ResetMDT", "Unlock"})
	t.AppendRow(table.Row{
		"WCC", l.boolean(wcc.Alarm), l.boolean(wcc.Reset), l.boolean(wcc.ResetMDT), l.boolean(wcc.Unlock),
	})
}

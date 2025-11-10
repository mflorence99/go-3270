package logger

import (
	"go3270/emulator/consts"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log WCC

func (l *Logger) logWCC(wcc consts.WCC) {
	t := l.newTable(text.FgHiBlue, "")
	defer t.Render()

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter},
		{Number: 4, Align: text.AlignCenter},
		{Number: 5, Align: text.AlignCenter},
	})

	t.AppendHeader(table.Row{
		"",
		"Alarm",
		"Reset",
		"ResetMDT",
		"Unlock",
	})

	t.AppendRow(table.Row{
		"WCC",
		l.boolean(wcc.Alarm),
		l.boolean(wcc.Reset),
		l.boolean(wcc.ResetMDT),
		l.boolean(wcc.Unlock),
	})
}

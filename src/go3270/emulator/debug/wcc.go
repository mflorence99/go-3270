package debug

import (
	"go3270/emulator/wcc"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func LogWCC(wcc *wcc.WCC) {
	t := NewTable()
	defer t.Render()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter},
		{Number: 4, Align: text.AlignCenter},
		{Number: 5, Align: text.AlignCenter},
	})
	t.AppendHeader(table.Row{"", "Alarm", "Reset", "ResetMDT", "Unlock"})
	t.AppendRows([]table.Row{
		{"WCC", Bool(wcc.Alarm), Bool(wcc.Reset), Bool(wcc.ResetMDT), Bool(wcc.Unlock)},
	})
}

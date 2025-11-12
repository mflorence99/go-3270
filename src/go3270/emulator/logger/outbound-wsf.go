package logger

import (
	"fmt"
	"go3270/emulator/sfld"
	"go3270/emulator/stream"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟦 Outbound WSF

func (l *Logger) logOutboundWSF(out *stream.Outbound, color text.Color) {
	t := l.newTable(color, "Outbound WSF")
	defer t.Render()

	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Info"})
	sflds := sfld.SFldsFromStream(out)
	for _, sfld := range sflds {
		t.AppendRow(table.Row{sfld.ID, fmt.Sprintf("% #x", sfld.Info)})
	}
}

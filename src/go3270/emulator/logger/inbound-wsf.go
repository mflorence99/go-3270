package logger

import (
	"bytes"
	"fmt"
	"go3270/emulator/sfld"
	"go3270/emulator/stream"
	"go3270/emulator/types"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟦 log WSF

// TODO 🔥 only really handles Query Reply

func (l *Logger) logInboundWSF(chars []byte) {
	t := l.newTable(text.FgHiGreen, ("Inbound WSF"))
	defer t.Render()

	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, types.LT)
	in := stream.NewOutbound(utils.Ternary(ok, slice, chars), l.bus)

	// 👇 eat the AID
	in.Next()

	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Type", "Info"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Transformer: l.wrap(60), WidthMax: 60, WidthMin: 80},
	})
	sflds := sfld.SFldsFromStream(in)

	for _, sfld := range sflds {
		switch {

		case sfld.ID == types.QUERY_REPLY:
			qcode := types.QCode(sfld.Info[0])
			t.AppendRow(table.Row{sfld.ID, qcode, fmt.Sprintf("% 02x", sfld.Info[1:])})

		default:
			t.AppendRow(table.Row{sfld.ID, "", fmt.Sprintf("% 02x", sfld.Info)})

		}
	}
}

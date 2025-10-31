package debug

import (
	"bytes"
	"fmt"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/stream"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logInbound(chars []byte) {
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, consts.LT)
	in := stream.NewOutbound(utils.Ternary(ok, slice, chars))
	char, _ := in.Next()
	aid := consts.AID(char)
	// 👇 short reads only contain the AID
	raw, ok := in.NextSlice(2)
	if ok {
		t := l.newTable(text.FgHiGreen, fmt.Sprintf("%s Inbound (3270 -> App)", aid))
		defer t.Render()
		// 👇 table rows
		t.AppendHeader(table.Row{"Row", "Col", "Data"})
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 3, Transformer: l.wrap(80), WidthMax: 80},
		})
		// 👇 one row just for the cursor
		cursorAt := conv.AddrFromBytes(raw)
		row, col := l.cfg.Addr2RC(cursorAt)
		t.AppendRow(table.Row{row, col, "(cursorAt)"})
		// 👇 we will aggregate data delimited by SBA's
		data := make([]byte, 0)
		// 👇 look at each byte to see if it is an order
		for in.HasNext() {
			char, _ := in.Next()
			order := consts.Order(char)
			switch order {

			case consts.SBA:
				if len(data) > 0 {
					t.AppendRow(table.Row{row, col, string(data)})
					data = make([]byte, 0)
				}
				raw, _ := in.NextSlice(2)
				addr := conv.AddrFromBytes(raw)
				row, col = l.cfg.Addr2RC(addr)

			default:
				data = append(data, conv.E2A(char))

			}
		}
		// 👇 don't forget the last field
		if len(data) > 0 {
			t.AppendRow(table.Row{row, col, string(data)})
		}
	} else {
		println(fmt.Sprintf("🐞 %s Short Read (3270 -> App)", aid))
	}
}

// TODO 🔥 only really handles Query Reply

func (l *Logger) logInboundWSF(chars []byte) {
	t := l.newTable(text.FgHiGreen, ("Inbound WSF (3270 -> App)"))
	defer t.Render()
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, consts.LT)
	in := stream.NewOutbound(utils.Ternary(ok, slice, chars))
	// 👇 eat the AID
	in.Next()
	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Type", "Info"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Transformer: l.wrap(60), WidthMax: 60},
	})
	sflds := consts.SFldsFromStream(in)
	for _, sfld := range sflds {
		switch {

		case sfld.ID == consts.QUERY_REPLY:
			qcode := consts.QCode(sfld.Info[0])
			t.AppendRow(table.Row{sfld.ID, qcode, fmt.Sprintf("% 02x", sfld.Info[1:])})

		default:
			t.AppendRow(table.Row{sfld.ID, "", fmt.Sprintf("% 02x", sfld.Info)})

		}
	}
}

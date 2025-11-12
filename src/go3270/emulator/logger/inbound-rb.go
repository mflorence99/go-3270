package logger

import (
	"bytes"
	"fmt"
	"go3270/emulator/conv"
	"go3270/emulator/stream"
	"go3270/emulator/types"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟦 log RB

func (l *Logger) logInboundRB(chars []byte) {
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, types.LT)
	in := stream.NewOutbound(utils.Ternary(ok, slice, chars), l.bus)
	char := in.MustNext()
	aid := types.AID(char)

	// 👇 create table
	t := l.newTable(text.FgHiGreen, fmt.Sprintf("%s Inbound RB", aid))
	defer t.Render()

	// 👇 table headers
	t.AppendHeader(table.Row{"", "Row", "Col", "Data"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, Transformer: l.wrap(80), WidthMax: 80, WidthMin: 80},
	})

	// 👇 one row just for the cursor
	raw := in.MustNextSlice(2)
	cursorAt := conv.Bytes2Addr(raw)
	row, col := l.cfg.Addr2RC(cursorAt)
	t.AppendRow(table.Row{"IC", row, col})

	// 👇 we will aggregate data delimited by SF and SFE's
	addr := 0
	data := make([]byte, 0)

	// 👇 look at each byte to see if it is an order
	var lastOrder types.Order
	for in.HasNext() {
		char := in.MustNext()
		order := types.Order(char)
		switch order {

		case types.SA:
			chars := in.MustNextSlice(2)
			attrs := types.NewExtendedAttrs(chars)
			row, col = l.cfg.Addr2RC(addr)
			t.AppendRow(table.Row{"SA", row, col, attrs.String()})

		case types.SF:
			if len(data) > 0 {
				t.AppendRow(table.Row{lastOrder.String(), row, col, string(data)})
				data = make([]byte, 0)
			}
			raw := in.MustNext()
			attrs := types.NewBasicAttrs(raw)
			row, col = l.cfg.Addr2RC(addr)
			lastOrder = order
			t.AppendRow(table.Row{lastOrder.String(), row, col, attrs.String()})

		case types.SFE:
			if len(data) > 0 {
				t.AppendRow(table.Row{lastOrder.String(), row, col, string(data)})
				data = make([]byte, 0)
			}
			count := in.MustNext()
			raw := in.MustNextSlice(int(count) * 2)
			attrs := types.NewExtendedAttrs(raw)
			row, col = l.cfg.Addr2RC(addr)
			lastOrder = order
			t.AppendRow(table.Row{lastOrder.String(), row, col, attrs.String()})

		default:
			data = append(data, conv.E2A(char))
			addr++

		}
	}

	// 👇 don't forget the last field
	if len(data) > 0 {
		t.AppendRow(table.Row{lastOrder.String(), row, col, string(data)})
	}

}

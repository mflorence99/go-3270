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

// 🟦 log RM/RMA

func (l *Logger) logInboundRM(chars []byte) {
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, types.LT)
	in := stream.NewOutbound(utils.Ternary(ok, slice, chars), l.bus)
	char := in.MustNext()
	aid := types.AID(char)

	// 👇 create table
	t := l.newTable(text.FgHiGreen, fmt.Sprintf("%s: %s Inbound RM/RMA", aid, l.buf.Mode()))
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

	// 👇 we will aggregate data delimited by SBA's
	addr := 0
	row, col = l.cfg.Addr2RC(addr)
	data := make([]byte, 0)

	// 👇 common code to print attributes
	appendAttrs := func(order types.Order, attrs *types.Attrs) {
		colorizer := text.Colors{text.FgYellow}
		row, col = l.cfg.Addr2RC(addr)
		t.AppendRow(table.Row{types.OrderFor(order), row, col, colorizer.Sprint(attrs.String())})
	}

	// 👇 common code to flush aggregated data
	flush := func(data []byte) []byte {
		if len(data) > 0 {
			t.AppendRow(table.Row{"", row, col, string(data)})
			return make([]byte, 0)
		}
		return data
	}

	// 👇 look at each byte to see if it is an order
	for in.HasNext() {
		char := in.MustNext()
		order := types.Order(char)
		switch order {

		case types.SA:
			data = flush(data)
			chars := in.MustNextSlice(2)
			attrs := types.NewExtendedAttrs(chars)
			appendAttrs(order, attrs)

		case types.SBA:
			data = flush(data)
			raw := in.MustNextSlice(2)
			addr = conv.Bytes2Addr(raw)
			row, col = l.cfg.Addr2RC(addr)
			t.AppendRow(table.Row{"SBA", row, col, ""})

		default:
			data = append(data, conv.E2A(char))
			addr++

		}
	}

	// 👇 don't forget the last field
	flush(data)
}

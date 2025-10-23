package debug

import (
	"fmt"
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/stream"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (l *Logger) logOutbound(chars []byte) {
	// 👇 analyze the commands in the stream
	out := stream.NewOutbound(chars)
	char, _ := out.Next()
	cmd := consts.Command(char)
	println(fmt.Sprintf("🐞 %s commanded", cmd))
	// 👇 if there's a hyte that follows, it's the WCC which we can ignore
	_, ok := out.Next()
	if ok {
		// 👇 now we can analyze commands with data
		switch cmd {
		case consts.WSF:
			l.logWSF(out)
		default:
			l.logOrders(out)
		}
	}
}

func (l *Logger) logOrders(out *stream.Outbound) {
	t := NewTable()
	defer t.Render()
	t.AppendHeader(table.Row{"Cmd", "Row", "Col", "Attributes"})
	fldAttrs := &attrs.Attrs{Protected: true}
	// 👇 look at each byte to see if it is an order
	for out.HasNext() {
		char, _ := out.Next()
		order := consts.Order(char)
		switch order {

		case consts.SBA:
			raw, _ := out.NextSlice(2)
			addr := conv.AddrFromBytes(raw)
			row, col := l.cfg.Addr2RC(addr)
			t.AppendRow(table.Row{order, fmt.Sprintf("%03d", row), fmt.Sprintf("%03d", col), ""})

		case consts.SF:
			next, _ := out.Next()
			fldAttrs = attrs.NewBasic(next)
			t.AppendRow(table.Row{order, "", "", fldAttrs})

		}
	}
}

func (l *Logger) logWSF(out *stream.Outbound) {}

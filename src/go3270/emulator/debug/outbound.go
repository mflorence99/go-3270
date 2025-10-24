package debug

import (
	"fmt"
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/stream"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logOutbound(chars []byte) {
	// 👇 analyze the commands in the stream
	out := stream.NewOutbound(chars)
	char, _ := out.Next()
	cmd := consts.Command(char)
	println(fmt.Sprintf("🐞 %s OUTBOUND", cmd))
	// 👇 now we can analyze commands with data
	switch cmd {

	case consts.EW:
		_, ok := out.Next() // 👈 eat thw WCC
		if ok {
			l.logOrders(out)
		}

	case consts.EWA:
		_, ok := out.Next() // 👈 eat thw WCC
		if ok {
			l.logOrders(out)
		}

	case consts.W:
		_, ok := out.Next() // 👈 eat thw WCC
		if ok {
			l.logOrders(out)
		}

	case consts.WSF:
		l.logWSF(out)

	default:
		l.logOrders(out)
	}
}

func (l *Logger) logOrders(out *stream.Outbound) {
	t := l.newTable(text.FgHiYellow)
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"Cmd", "Row", "Col", "SF", "Blink", "Color", "Hidden", "Hilite", "MDT", "Num", "Prot", "Rev", "Uscore"})
	addr := 0
	a := &attrs.Attrs{Protected: true}
	// 👇 look at each byte to see if it is an order
	for out.HasNext() {
		char, _ := out.Next()
		order := consts.Order(char)
		switch order {

		case consts.EUA:

		case consts.GE:

		case consts.IC:
			l.withoutAttrs(t, order, addr)

		case consts.MF:

		case consts.PT:

		case consts.RA:
			l.withoutAttrs(t, order, addr)
			raw, _ := out.NextSlice(2)
			addr = conv.AddrFromBytes(raw)

		case consts.SA:
			bytes, _ := out.NextSlice(2)
			a = attrs.NewModified(a, bytes)
			l.withAttrs(t, order, addr, a, false)

		case consts.SBA:
			raw, _ := out.NextSlice(2)
			addr = conv.AddrFromBytes(raw)

		case consts.SF:
			next, _ := out.Next()
			a = attrs.NewBasic(next)
			l.withAttrs(t, order, addr, a, true)
			addr++

		case consts.SFE:
			count, _ := out.Next()
			next, _ := out.NextSlice(int(count) * 2)
			a = attrs.NewExtended(next)
			l.withAttrs(t, order, addr, a, true)
			addr++

		default:
			addr++

		}
	}
}

func (l *Logger) logWSF(out *stream.Outbound) {
	t := l.newTable(text.FgHiYellow)
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Info"})
	sflds := consts.FromStream(out)
	for _, sfld := range sflds {
		t.AppendRow(table.Row{sfld.ID, fmt.Sprintf("% #v", sfld.Info)})
	}
}

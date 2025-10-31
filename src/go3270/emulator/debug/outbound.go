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
	// 👇 now we can analyze commands with data
	switch cmd {

	case consts.EW:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOrders(out, cmd)
		}

	case consts.EWA:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOrders(out, cmd)
		}

	case consts.W:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOrders(out, cmd)
		}

	case consts.WSF:
		l.logOutboundWSF(out)
	}
}

func (l *Logger) logOrders(out *stream.Outbound, cmd consts.Command) {
	t := l.newTable(text.FgHiYellow, fmt.Sprintf("%s Outbound (App -> 3270)\nNOTE: RA orders are listed in pairs, one for start, the second for stop", cmd))
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"Order", "Row", "Col", "SF", "Blink", "Color", "Hidden", "Hilite", "MDT", "Num", "Prot", "Rev", "Uscore"})
	addr := 0
	fldAttrs := &attrs.Attrs{Protected: true}
	// 👇 look at each byte to see if it is an order
	for out.HasNext() {
		char, _ := out.Next()
		order := consts.Order(char)
		switch order {

		case consts.EUA:

		case consts.GE:

		case consts.IC:
			l.withoutAttrs(t, order, addr, ' ')

		case consts.MF:

		case consts.PT:

		case consts.RA:
			raw, _ := out.NextSlice(2)
			char, _ := out.Next()
			l.withoutAttrs(t, order, addr, char)
			addr = conv.AddrFromBytes(raw)
			l.withoutAttrs(t, order, addr, char)

		case consts.SA:
			bytes, _ := out.NextSlice(2)
			fldAttrs = attrs.NewModified(fldAttrs, bytes)
			l.withAttrs(t, order, addr, fldAttrs, false)

		case consts.SBA:
			raw, _ := out.NextSlice(2)
			addr = conv.AddrFromBytes(raw)
			l.withoutAttrs(t, order, addr, 0)

		case consts.SF:
			next, _ := out.Next()
			fldAttrs = attrs.NewBasic(next)
			l.withAttrs(t, order, addr, fldAttrs, true)
			addr++

		case consts.SFE:
			count, _ := out.Next()
			next, _ := out.NextSlice(int(count) * 2)
			fldAttrs = attrs.NewExtended(next)
			l.withAttrs(t, order, addr, fldAttrs, true)
			addr++

		default:
			addr++

		}
	}
}

func (l *Logger) logOutboundWSF(out *stream.Outbound) {
	t := l.newTable(text.FgHiYellow, "Outbound WSF (App -> 3270)")
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Info"})
	sflds := consts.SFldsFromStream(out)
	for _, sfld := range sflds {
		t.AppendRow(table.Row{sfld.ID, fmt.Sprintf("% #v", sfld.Info)})
	}
}

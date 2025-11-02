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

// 🟧 Debugger: log outbound (3270 <- app) stream

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
			l.logOutboundOrders(out, cmd)
		}

	case consts.EWA:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd)
		}

	case consts.W:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd)
		}

	case consts.WSF:
		l.logOutboundWSF(out)
	}
}

func (l *Logger) logOutboundOrders(out *stream.Outbound, cmd consts.Command) {
	t := l.newTable(text.FgHiYellow, fmt.Sprintf("%s Outbound (App -> 3270)\nNOTE: EUA and RA orders are listed in start/stop pairs", cmd))
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
			raw, _ := out.NextSlice(2)
			l.withoutAttrs(t, order, addr, ' ')
			addr = conv.AddrFromBytes(raw)
			l.withoutAttrs(t, order, addr, ' ')

		case consts.GE:
			l.withoutAttrs(t, order, addr, ' ')

		case consts.IC:
			l.withoutAttrs(t, order, addr, ' ')

		case consts.MF:
			count, _ := out.Next()
			raw, _ := out.NextSlice(int(count) * 2)
			fldAttrs = attrs.NewExtended(raw)
			l.withAttrs(t, order, addr, fldAttrs, true)
			addr++

		case consts.PT:
			l.withoutAttrs(t, order, addr, ' ')

		case consts.RA:
			raw, _ := out.NextSlice(2)
			char, _ := out.Next()
			l.withoutAttrs(t, order, addr, char)
			addr = conv.AddrFromBytes(raw)
			l.withoutAttrs(t, order, addr, char)

		case consts.SA:
			chars, _ := out.NextSlice(2)
			fldAttrs = attrs.NewModified(fldAttrs, chars)
			l.withAttrs(t, order, addr, fldAttrs, false)

		case consts.SBA:
			raw, _ := out.NextSlice(2)
			addr = conv.AddrFromBytes(raw)
			l.withoutAttrs(t, order, addr, 0)

		case consts.SF:
			raw, _ := out.Next()
			fldAttrs = attrs.NewBasic(raw)
			l.withAttrs(t, order, addr, fldAttrs, true)
			addr++

		case consts.SFE:
			count, _ := out.Next()
			raw, _ := out.NextSlice(int(count) * 2)
			fldAttrs = attrs.NewExtended(raw)
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

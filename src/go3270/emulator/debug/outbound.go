package debug

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/stream"
	"go3270/emulator/utils"

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
			l.logOutboundOrders(out, cmd, text.FgHiYellow)
		}

	case consts.EWA:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgHiYellow)
		}

	case consts.W:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case consts.WSF:
		l.logOutboundWSF(out, text.FgCyan)
	}
}

func (l *Logger) logOutboundOrders(out *stream.Outbound, cmd consts.Command, color text.Color) {
	t := l.newTable(color, fmt.Sprintf("%s Outbound (App -> 3270)\nNOTE: EUA and RA orders are listed in start/stop pairs", cmd))
	defer t.Render()
	addr := 0
	fldAddr := 0
	fldAttrs := &consts.Attrs{Protected: true}
	// 👇 header
	t.AppendHeader(table.Row{
		"",
		"Row",
		"Col",
		"SF",
		"Blink",
		"Color",
		"Hidden",
		"Hilite",
		"MDT",
		"Num",
		"Prot",
		"Rev",
		"Uscore",
		"Out",
		"LCID",
	})
	// 👇 look at each byte to see if it is an order
	for out.HasNext() {
		char, _ := out.Next()
		order := consts.Order(char)
		switch order {

		case consts.EUA:
			raw, _ := out.NextSlice(2)
			l.withoutAttrs(t, order, addr, ' ')
			addr = conv.Bytes2Addr(raw)
			l.withoutAttrs(t, order, addr, ' ')

		case consts.GE:
			char, _ := out.Next()
			l.withoutAttrs(t, order, addr, conv.E2A(char))

		case consts.IC:
			l.withoutAttrs(t, order, addr, ' ')

		case consts.MF:
			count, _ := out.Next()
			raw, _ := out.NextSlice(int(count) * 2)
			fldAttrs = consts.NewExtendedAttrs(raw)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr}
			l.withAttrs(t, order, addr, cell)
			addr++

		case consts.PT:
			l.withoutAttrs(t, order, addr, ' ')

		case consts.RA:
			raw, _ := out.NextSlice(2)
			char, _ := out.Next()
			if consts.Order(char) == consts.GE {
				char, _ = out.Next()
				l.withoutAttrs(t, consts.GE, addr, conv.E2A(char))
			}
			l.withoutAttrs(t, order, addr, conv.E2A(char))
			addr = conv.Bytes2Addr(raw)
			l.withoutAttrs(t, order, addr, conv.E2A(char))

		case consts.SA:
			chars, _ := out.NextSlice(2)
			fldAttrs = consts.NewModifiedAttrs(fldAttrs, chars)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr}
			l.withAttrs(t, order, addr, cell)

		case consts.SBA:
			raw, _ := out.NextSlice(2)
			addr = conv.Bytes2Addr(raw)
			l.withoutAttrs(t, order, addr, 0)

		case consts.SF:
			raw, _ := out.Next()
			fldAddr = addr
			fldAttrs = consts.NewBasicAttrs(raw)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr, FldStart: true}
			l.withAttrs(t, order, addr, cell)
			addr++

		case consts.SFE:
			count, _ := out.Next()
			raw, _ := out.NextSlice(int(count) * 2)
			fldAddr = addr
			fldAttrs = consts.NewExtendedAttrs(raw)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr, FldStart: true}
			l.withAttrs(t, order, addr, cell)
			addr++

		default:
			addr++

		}
	}
}

func (l *Logger) logOutboundWSF(out *stream.Outbound, color text.Color) {
	t := l.newTable(color, "Outbound WSF (App -> 3270)")
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Info"})
	sflds := consts.SFldsFromStream(out)
	for _, sfld := range sflds {
		t.AppendRow(table.Row{sfld.ID, fmt.Sprintf("% #x", sfld.Info)})
	}
}

func (l *Logger) withAttrs(t table.Writer, cmd any, addr int, cell *buffer.Cell) {
	row, col := l.cfg.Addr2RC(addr)
	t.AppendRow(table.Row{
		cmd,
		row,
		col,
		l.boolean(cell.FldStart),
		utils.Ternary(cell.Attrs.Blink, "BLINK", ""),
		utils.Ternary(cell.Attrs.Color != 0x00, consts.ColorFor(cell.Attrs.Color), ""),
		utils.Ternary(cell.Attrs.Hidden, "HIDDEN", ""),
		utils.Ternary(cell.Attrs.Highlight, "HILITE", ""),
		utils.Ternary(cell.Attrs.Modified, "MDT", ""),
		utils.Ternary(cell.Attrs.Numeric, "NUM", ""),
		utils.Ternary(cell.Attrs.Protected, "PROT", ""),
		utils.Ternary(cell.Attrs.Reverse, "REV", ""),
		utils.Ternary(cell.Attrs.Underscore, "USCORE", ""),
		utils.Ternary(cell.Attrs.Outline != 0x00, consts.OutlineFor(cell.Attrs.Outline), ""),
		utils.Ternary(cell.Attrs.LCID != 0x00, cell.Attrs.LCID.String(), ""),
	})
}

func (l *Logger) withoutAttrs(t table.Writer, cmd any, addr int, char byte) {
	row, col := l.cfg.Addr2RC(addr)
	t.AppendRow(table.Row{
		cmd,
		row,
		col,
		utils.Ternary(char >= 0x20, string(char), " "),
	})
}

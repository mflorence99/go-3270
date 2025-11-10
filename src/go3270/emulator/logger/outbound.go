package logger

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/conv"
	"go3270/emulator/sfld"
	"go3270/emulator/stream"
	"go3270/emulator/types"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log outbound (3270 <- app) stream

func (l *Logger) logOutbound(chars []byte) {
	// 👇 analyze the commands in the stream
	out := stream.NewOutbound(chars, l.bus)
	char := out.MustNext()
	cmd := types.Command(char)
	// 👇 now we can analyze commands with data
	switch cmd {

	case types.EW:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.EWA:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.W:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.WSF:
		l.logOutboundWSF(out, text.FgCyan)
	}
}

func (l *Logger) logOutboundOrders(out *stream.Outbound, cmd types.Command, color text.Color) {
	t := l.newTable(color, fmt.Sprintf("%s Outbound (App -> 3270)\nNOTE: EUA and RA orders are listed in start/stop pairs", cmd))
	defer t.Render()
	addr := 0
	fldAddr := 0
	fldAttrs := &types.Attrs{Protected: true}

	// 👇 header
	t.AppendHeader(table.Row{
		"",
		"Row",
		"Col",
		"SF",
		"Skip",
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
		char := out.MustNext()
		order := types.Order(char)
		switch order {

		case types.EUA:
			raw := out.MustNextSlice(2)
			l.withoutAttrs(t, order, addr, ' ')
			addr = conv.Bytes2Addr(raw)
			l.withoutAttrs(t, order, addr, ' ')

		case types.GE:
			char := out.MustNext()
			l.withoutAttrs(t, order, addr, conv.E2A(char))

		case types.IC:
			l.withoutAttrs(t, order, addr, ' ')

		case types.MF:
			count := out.MustNext()
			raw := out.MustNextSlice(int(count) * 2)
			fldAttrs = types.NewExtendedAttrs(raw)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr}
			l.withAttrs(t, order, addr, cell)
			addr++

		case types.PT:
			l.withoutAttrs(t, order, addr, ' ')

		case types.RA:
			raw := out.MustNextSlice(2)
			char := out.MustNext()
			if types.Order(char) == types.GE {
				char = out.MustNext()
				l.withoutAttrs(t, types.GE, addr, conv.E2A(char))
			}
			l.withoutAttrs(t, order, addr, conv.E2A(char))
			addr = conv.Bytes2Addr(raw)
			l.withoutAttrs(t, order, addr, conv.E2A(char))

		case types.SA:
			chars := out.MustNextSlice(2)
			fldAttrs = types.NewModifiedAttrs(fldAttrs, chars)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr}
			l.withAttrs(t, order, addr, cell)

		case types.SBA:
			raw := out.MustNextSlice(2)
			addr = conv.Bytes2Addr(raw)
			l.withoutAttrs(t, order, addr, 0)

		case types.SF:
			raw := out.MustNext()
			fldAddr = addr
			fldAttrs = types.NewBasicAttrs(raw)
			cell := &buffer.Cell{Attrs: fldAttrs, FldAddr: fldAddr, FldStart: true}
			l.withAttrs(t, order, addr, cell)
			addr++

		case types.SFE:
			count := out.MustNext()
			raw := out.MustNextSlice(int(count) * 2)
			fldAddr = addr
			fldAttrs = types.NewExtendedAttrs(raw)
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
	sflds := sfld.SFldsFromStream(out)
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
		utils.Ternary(cell.Attrs.Autoskip, "SKIP", ""),
		utils.Ternary(cell.Attrs.Blink, "BLINK", ""),
		utils.Ternary(cell.Attrs.Color != 0x00, types.ColorFor(cell.Attrs.Color), ""),
		utils.Ternary(cell.Attrs.Hidden, "HIDDEN", ""),
		utils.Ternary(cell.Attrs.Highlight, "HILITE", ""),
		utils.Ternary(cell.Attrs.MDT, "MDT", ""),
		utils.Ternary(cell.Attrs.Numeric, "NUM", ""),
		utils.Ternary(cell.Attrs.Protected, "PROT", ""),
		utils.Ternary(cell.Attrs.Reverse, "REV", ""),
		utils.Ternary(cell.Attrs.Underscore, "USCORE", ""),
		utils.Ternary(cell.Attrs.Outline != 0x00, types.OutlineFor(cell.Attrs.Outline), ""),
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

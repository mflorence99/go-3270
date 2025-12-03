package core

import (
	"bytes"
	"emulator/conv"
	"emulator/types"
	"emulator/utils"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Subscribe to important messages and log debugging data

// TODO 🔥 I never wanted all this in one file, but putting it in its
// own package causes a cyclic dependency only a REALLY ugly
// compromise removes. So this is the least-bad (so far) choice.

// 🔥 Most logging avoids blocking the main thread by using Go routines

type Logger struct {
	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewLogger(emu *Emulator) *Logger {
	l := new(Logger)
	l.emu = emu
	// 👇 subscriptions
	l.emu.Bus.SubClose(l.close)
	l.emu.Bus.SubInitialize(l.initialize)
	l.emu.Bus.SubInbound(l.inbound)
	l.emu.Bus.SubOutbound(l.outbound)
	l.emu.Bus.SubProbe(l.probe)
	l.emu.Bus.SubRender(l.render)
	l.emu.Bus.SubTrace(l.trace)
	l.emu.Bus.SubWCChar(l.wcc)
	return l
}

func (l *Logger) close() {
	println("🐞 Emulator closed")
}

func (l *Logger) initialize() {
	println("🐞 Emulator initialized")
	go func() {
		l.logConfig()
		l.logCLUT()
	}()
}

func (l *Logger) inbound(chars []byte, hints PubInboundHints) {
	go func() {
		// 👇 supplement with an old-fashioned core dump
		l.dump(text.FgHiGreen, "Inbound Core Dump", chars)
		l.logInbound(chars, hints)
	}()
}

func (l *Logger) outbound(chars []byte) {
	go func() {
		// 👇 supplement with an old-fashioned core dump
		l.dump(text.FgYellow, "Outbound Core Dump", chars)
		l.logOutbound(chars)
	}()
}

func (l *Logger) probe(addr uint) {
	go func() {
		l.logProbe(addr)
	}()
}

func (l *Logger) render() {
	go func() {
		l.logBuffer()
		l.logFlds()
	}()
}

func (l *Logger) trace(topic Topic, handler interface{}) {
	// 👇 we need this to be synchronous to make sense
	l.logTrace(topic, handler)
}

func (l *Logger) wcc(wcc types.WCC) {
	go func() {
		l.logWCC(wcc)
	}()
}

// ---------------------------------------------------------------------------
// 🟦 Dump a slice of bytes
// ---------------------------------------------------------------------------

func (l *Logger) dump(color text.Color, title string, chars []byte) {
	t := l.newTable(color, title)
	defer t.Render()

	// 👇 Header rows

	t.AppendHeader(table.Row{
		"",
		"0        0        0        0        1        1        1        1 ",
		"0                1",
	})

	t.AppendHeader(table.Row{
		"",
		"0 1 2 3  4 5 6 7  8 9 a b  c d e f  0 1 2 3  4 5 6 7  8 9 a b  c d e f",
		"0123456789abcdef 0123456789abcdef",
	})

	// 👇 Dump rows in 32-byte wide slices

	for offset := 0; offset < len(chars); offset += 32 {
		slice := chars[offset:min(offset+32, len(chars))]

		var hex strings.Builder
		for ix := 0; ix < len(slice); ix += 4 {
			chunk := slice[ix:min(ix+4, len(slice))]
			hex.WriteString(fmt.Sprintf("%x ", chunk))
		}

		var ascii strings.Builder
		for ix := 0; ix < len(slice); ix += 16 {
			chunk := slice[ix:min(ix+16, len(slice))]
			// 🔥 we don't know the LCID here, so just dump as ASCII
			ascii.WriteString(fmt.Sprintf("%s ", conv.E2As(string(chunk))))
		}

		t.AppendRow(table.Row{
			fmt.Sprintf("%06x", offset),
			hex.String(),
			ascii.String(),
		})
	}
}

// ---------------------------------------------------------------------------
// 🟦 Log buffer contents
// ---------------------------------------------------------------------------

func (l *Logger) logBuffer() {

	cursor := fmt.Sprintf("%s%s", text.FgHiGreen.Sprint("\u25a0"), text.FgWhite.Sprint("\u200b"))

	protected := fmt.Sprintf("%s%s", text.FgHiCyan.Sprint("\u00b6"), text.FgWhite.Sprint("\u200b"))

	unprotected := fmt.Sprintf("%s%s", text.FgHiRed.Sprint("\u00b6"), text.FgWhite.Sprint("\u200b"))

	// 👇 define the table
	t := l.newTable(text.FgHiBlue, fmt.Sprintf("%s Buffer\nCURSOR: %s PROT: %s UNPROT: %s", l.emu.Buf.Mode(), cursor, protected, unprotected))
	defer t.Render()

	// 👇 header rows
	row1 := ""
	row2 := ""
	for ix := uint(10); ix <= l.emu.Cfg.Cols; ix += 10 {
		row1 += fmt.Sprintf("%10d", ix/10)
		row2 += "1234567890"
	}
	t.AppendHeader(table.Row{
		"",
		fmt.Sprintf("%s\n%s", row1, row2),
	})

	// 👇 where's the cursorAt?
	row, col := l.emu.Cfg.Addr2RC(l.emu.State.Status.CursorAt)

	// 👇 data rows
	for iy := uint(1); iy <= l.emu.Cfg.Rows; iy++ {
		var b strings.Builder
		// 👇 data cols
		for ix := uint(1); ix <= l.emu.Cfg.Cols; ix++ {
			// 👇 show the cursor specially
			if iy == row && ix == col {
				b.WriteString(cursor)
			} else {
				// 👇 or the cell contents, best as we can
				cell, ok := l.emu.Buf.Peek(ix + ((iy - 1) * l.emu.Cfg.Cols) - 1)
				if cell != nil && ok {

					if cell.IsFldStart() {
						b.WriteString(utils.Ternary(cell.Attrs.Protected, protected, unprotected))

					} else {
						str := " "
						if cell.Char > 0x40 {
							str = string(conv.E2Rune(cell.Attrs.LCID, cell.Char))
						}
						if cell.Attrs.CharAttr {
							str = fmt.Sprintf("%s%s", text.FgYellow.Sprint(str), text.FgWhite.Sprint("\u200b"))
						}
						b.WriteString(str)
					}
				}
			}
		}
		t.AppendRow(table.Row{iy, b.String()})
	}
}

// ---------------------------------------------------------------------------
// 🟦 Log configuration
// ---------------------------------------------------------------------------

func (l *Logger) logConfig() {
	t := l.newTable(text.FgHiBlue, "Config")
	defer t.Render()

	t.AppendHeader(table.Row{
		"",
		"",
		"",
		"BG",
		"",
		"Font",
		"Font",
		"Font",
		"Padding",
		"Padding",
		"",
	})
	t.AppendHeader(table.Row{
		"",
		"#Rows",
		"#Cols",
		"Color",
		"Mono",
		"Width",
		"Height",
		"Size",
		"Width",
		"Height",
		"Test",
	})

	t.AppendRow(table.Row{
		"CFG",
		l.emu.Cfg.Rows,
		l.emu.Cfg.Cols,
		l.emu.Cfg.BgColor,
		l.boolean(l.emu.Cfg.Monochrome),
		l.emu.Cfg.FontWidth,
		l.emu.Cfg.FontHeight,
		l.emu.Cfg.FontSize,
		l.emu.Cfg.PaddedWidth,
		l.emu.Cfg.PaddedHeight,
		l.emu.Cfg.Testpage,
	},
	)
}

func (l *Logger) logCLUT() {
	t := l.newTable(text.FgHiBlue, "CLUT")
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "Attr", "Color"})
	for k, v := range l.emu.Cfg.CLUT {
		t.AppendRow(table.Row{k, fmt.Sprintf("%#02x", byte(k)), v})
	}
}

// ---------------------------------------------------------------------------
// 🟦 Log buffer fields
// ---------------------------------------------------------------------------

func (l *Logger) logFlds() {
	t := l.newTable(text.FgHiBlue, "Buffer Fields")
	defer t.Render()

	// 👇 header rows
	t.AppendHeader(table.Row{
		"Row",
		"Col",
		"Len",
		"HIDDEN",
		"MDT",
		"PROT",
		"Data",
	})

	// 👇 data rows
	for _, fld := range l.emu.Flds.Flds {
		sf := fld.Cells[0]
		addr, _ := sf.GetFldAddr()
		row, col := l.emu.Cfg.Addr2RC(addr)
		// 👇 gather all the chars in the fld
		t.AppendRow(table.Row{
			row,
			col,
			len(fld.Cells),
			utils.Ternary(sf.Attrs.Hidden, "HIDDEN", ""),
			utils.Ternary(sf.Attrs.MDT, "MDT", ""),
			utils.Ternary(sf.Attrs.Protected, "PROT", ""),
			utils.Truncate(fld.String(), 60),
		})
	}
}

// ---------------------------------------------------------------------------
// 🟦 Log inbound data
// ---------------------------------------------------------------------------

func (l *Logger) logInbound(chars []byte, hints PubInboundHints) {
	switch {

	case hints.RB:
		l.logInboundRB(chars)

	case hints.RM:
		l.logInboundRM(chars)

	case hints.Short:
		l.logInboundShort(chars)

	case hints.WSF:
		l.logInboundWSF(chars)

	}
}

// ---------------------------------------------------------------------------
// 🟪 ...RB (read buffer)
// ---------------------------------------------------------------------------

func (l *Logger) logInboundRB(chars []byte) {
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, types.LT)
	in := NewOutbound(utils.Ternary(ok, slice, chars), l.emu.Bus)
	char := in.MustNext()
	aid := types.AID(char)

	// 👇 create table
	t := l.newTable(text.FgHiGreen, fmt.Sprintf("%s: %s Inbound RB", aid, l.emu.Buf.Mode()))
	defer t.Render()

	// 👇 table headers
	t.AppendHeader(table.Row{"", "Row", "Col", "Data"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, Transformer: l.wrap(80), WidthMax: 80, WidthMin: 80},
	})

	// 👇 one row just for the cursor
	raw := in.MustNextSlice(2)
	cursorAt := conv.Bytes2Addr(raw)
	row, col := l.emu.Cfg.Addr2RC(cursorAt)
	t.AppendRow(table.Row{"IC", row, col})

	// 👇 we will aggregate data delimited by SF and SFE's
	var addr uint
	row, col = l.emu.Cfg.Addr2RC(addr)
	data := make([]byte, 0)

	// 👇 common code to print attributes
	appendAttrs := func(order types.Order, attrs *types.Attrs) {
		colorizer := text.Colors{text.FgYellow}
		row, col = l.emu.Cfg.Addr2RC(addr)
		t.AppendRow(table.Row{types.OrderFor(order), row, col, colorizer.Sprint(attrs.String())})
		if order != types.SA {
			addr++
			row, col = l.emu.Cfg.Addr2RC(addr)
		}
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

		case types.SF:
			data = flush(data)
			raw := in.MustNext()
			attrs := types.NewBasicAttrs(raw)
			appendAttrs(order, attrs)

		case types.SFE:
			data = flush(data)
			count := in.MustNext()
			raw := in.MustNextSlice(int(count) * 2)
			attrs := types.NewExtendedAttrs(raw)
			appendAttrs(order, attrs)

		default:
			data = append(data, conv.E2A(char))
			addr++

		}
	}

	// 👇 don't forget the last field
	flush(data)

}

// ---------------------------------------------------------------------------
// 🟪 ...RM (read modified)
// ---------------------------------------------------------------------------

func (l *Logger) logInboundRM(chars []byte) {
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, types.LT)
	in := NewOutbound(utils.Ternary(ok, slice, chars), l.emu.Bus)
	char := in.MustNext()
	aid := types.AID(char)

	// 👇 create table
	t := l.newTable(text.FgHiGreen, fmt.Sprintf("%s: %s Inbound RM/RMA", aid, l.emu.Buf.Mode()))
	defer t.Render()

	// 👇 table headers
	t.AppendHeader(table.Row{"", "Row", "Col", "Data"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, Transformer: l.wrap(80), WidthMax: 80, WidthMin: 80},
	})

	// 👇 one row just for the cursor
	raw := in.MustNextSlice(2)
	cursorAt := conv.Bytes2Addr(raw)
	row, col := l.emu.Cfg.Addr2RC(cursorAt)
	t.AppendRow(table.Row{"IC", row, col})

	// 👇 we will aggregate data delimited by SBA's
	var addr uint
	row, col = l.emu.Cfg.Addr2RC(addr)
	data := make([]byte, 0)

	// 👇 common code to print attributes
	appendAttrs := func(order types.Order, attrs *types.Attrs) {
		colorizer := text.Colors{text.FgYellow}
		row, col = l.emu.Cfg.Addr2RC(addr)
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
			row, col = l.emu.Cfg.Addr2RC(addr)
			t.AppendRow(table.Row{"SBA", row, col, ""})

		default:
			data = append(data, conv.E2A(char))
			addr++

		}
	}

	// 👇 don't forget the last field
	flush(data)
}

// ---------------------------------------------------------------------------
// 🟪 ...short read
// ---------------------------------------------------------------------------

func (l *Logger) logInboundShort(chars []byte) {
	aid := types.AID(chars[0])
	fmt.Printf("🐞 %s Short Read\n", aid)
}

// ---------------------------------------------------------------------------
// 🟪 ...WSF
//
// TODO 🔥 only really handles Query Reply
// ---------------------------------------------------------------------------

func (l *Logger) logInboundWSF(chars []byte) {
	t := l.newTable(text.FgHiGreen, ("Inbound WSF"))
	defer t.Render()

	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, types.LT)
	in := NewOutbound(utils.Ternary(ok, slice, chars), l.emu.Bus)

	// 👇 eat the AID
	in.Next()

	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Type", "Info"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Transformer: l.wrap(60), WidthMax: 60, WidthMin: 80},
	})
	sflds := SFldsFromStream(in)

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

// ---------------------------------------------------------------------------
// 🟦 Log outbound data
// ---------------------------------------------------------------------------

func (l *Logger) logOutbound(chars []byte) {
	// 👇 analyze the commands in the stream
	out := NewOutbound(chars, l.emu.Bus)
	char := out.MustNext()
	cmd := types.Command(char)
	// 👇 now we can analyze commands with data
	switch cmd {

	case types.EW:
		if _, ok := out.Next(); ok { // 👈 eat the WCC
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.EWA:
		if _, ok := out.Next(); ok { // 👈 eat the WCC
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.W:
		if _, ok := out.Next(); ok { // 👈 eat the WCC
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.WSF:
		l.logOutboundWSF(out, text.FgCyan)
	}
}

// ---------------------------------------------------------------------------
// 🟪 ...orders
// ---------------------------------------------------------------------------

func (l *Logger) logOutboundOrders(out *Outbound, cmd types.Command, color text.Color) {
	t := l.newTable(color, fmt.Sprintf("%s Outbound Orders\nNOTE: EUA and RA orders are listed in start/stop pairs", cmd))
	defer t.Render()
	var addr uint
	fldAttrs := types.NewDefaultAttrs()

	// 👇 header
	t.AppendHeader(table.Row{
		"",
		"Row",
		"Col",
		"SF",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
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
			l.logOutboundOrdersWithoutAttrs(t, order, addr, ' ')
			addr = conv.Bytes2Addr(raw)
			l.logOutboundOrdersWithoutAttrs(t, order, addr, ' ')

		case types.GE:
			char := out.MustNext()
			l.logOutboundOrdersWithoutAttrs(t, order, addr, conv.E2A(char))

		case types.IC:
			l.logOutboundOrdersWithoutAttrs(t, order, addr, ' ')

		case types.MF:
			count := out.MustNext()
			raw := out.MustNextSlice(int(count) * 2)
			fldAttrs = types.NewExtendedAttrs(raw)
			l.logOutboundOrdersWithAttrs(t, order, addr, fldAttrs, false)
			addr++

		case types.PT:
			l.logOutboundOrdersWithoutAttrs(t, order, addr, ' ')

		case types.RA:
			raw := out.MustNextSlice(2)
			char := out.MustNext()
			if types.Order(char) == types.GE {
				char = out.MustNext()
				l.logOutboundOrdersWithoutAttrs(t, types.GE, addr, conv.E2A(char))
			}
			l.logOutboundOrdersWithoutAttrs(t, order, addr, conv.E2A(char))
			addr = conv.Bytes2Addr(raw)
			l.logOutboundOrdersWithoutAttrs(t, order, addr, conv.E2A(char))

		case types.SA:
			chars := out.MustNextSlice(2)
			fldAttrs = types.NewModifiedAttrs(fldAttrs, chars)
			l.logOutboundOrdersWithAttrs(t, order, addr, fldAttrs, false)

		case types.SBA:
			raw := out.MustNextSlice(2)
			addr = conv.Bytes2Addr(raw)
			l.logOutboundOrdersWithoutAttrs(t, order, addr, 0)

		case types.SF:
			raw := out.MustNext()
			fldAttrs = types.NewBasicAttrs(raw)
			l.logOutboundOrdersWithAttrs(t, order, addr, fldAttrs, true)
			addr++

		case types.SFE:
			count := out.MustNext()
			raw := out.MustNextSlice(int(count) * 2)
			fldAttrs = types.NewExtendedAttrs(raw)
			l.logOutboundOrdersWithAttrs(t, order, addr, fldAttrs, true)
			addr++

		default:
			addr++

		}
	}
}

func (l *Logger) logOutboundOrdersWithAttrs(t table.Writer, cmd any, addr uint, fldAttrs *types.Attrs, fldStart bool) {
	row, col := l.emu.Cfg.Addr2RC(addr)
	t.AppendRow(table.Row{
		cmd,
		row,
		col,
		l.boolean(fldStart),
		utils.Ternary(fldAttrs.Autoskip, "SKIP", ""),
		utils.Ternary(fldAttrs.Blink, "BLINK", ""),
		utils.Ternary(fldAttrs.Color != 0x00, types.ColorFor(fldAttrs.Color), ""),
		utils.Ternary(fldAttrs.Hidden, "HIDDEN", ""),
		utils.Ternary(fldAttrs.Highlight, "HILITE", ""),
		utils.Ternary(fldAttrs.Intensify, "INTENSE", ""),
		utils.Ternary(fldAttrs.MDT, "MDT", ""),
		utils.Ternary(fldAttrs.Numeric, "NUM", ""),
		utils.Ternary(fldAttrs.Protected, "PROT", ""),
		utils.Ternary(fldAttrs.Reverse, "REV", ""),
		utils.Ternary(fldAttrs.Underscore, "USCORE", ""),
		utils.Ternary(fldAttrs.Outline != 0x00, types.OutlineFor(fldAttrs.Outline), ""),
		utils.Ternary(fldAttrs.LCID != 0x00, fldAttrs.LCID.String(), ""),
	})
}

func (l *Logger) logOutboundOrdersWithoutAttrs(t table.Writer, cmd any, addr uint, char byte) {
	row, col := l.emu.Cfg.Addr2RC(addr)
	t.AppendRow(table.Row{
		cmd,
		row,
		col,
		utils.Ternary(char >= 0x20, string(char), " "),
	})
}

// ---------------------------------------------------------------------------
// 🟪 ...WSF
// ---------------------------------------------------------------------------

func (l *Logger) logOutboundWSF(out *Outbound, color text.Color) {
	t := l.newTable(color, "Outbound WSF")
	defer t.Render()

	// 👇 table rows
	t.AppendHeader(table.Row{"ID", "Info"})
	sflds := SFldsFromStream(out)
	for _, sfld := range sflds {
		t.AppendRow(table.Row{sfld.ID, fmt.Sprintf("% #x", sfld.Info)})
	}
}

// ---------------------------------------------------------------------------
// 🟦 Probe cell contents
// ---------------------------------------------------------------------------

func (l *Logger) logProbe(addr uint) {
	cell := l.emu.Buf.MustPeek(addr)
	t := l.newTable(utils.Ternary(cell.Attrs.CharAttr, text.FgRed, text.FgHiRed), "")
	defer t.Render()

	// 👇 header
	t.AppendHeader(table.Row{
		"",
		"",
		"Fld",
		"Fld",
		"Cell",
		"Cell",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	})
	t.AppendHeader(table.Row{
		"",
		"SF",
		"Row",
		"Col",
		"Row",
		"Col",
		"SA",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"Out",
		"LCID",
	})

	// 👇 extract data
	crow, ccol := l.emu.Cfg.Addr2RC(addr)
	fldAddr, _ := cell.GetFldAddr()
	frow, fcol := l.emu.Cfg.Addr2RC(fldAddr)
	char := fmt.Sprintf("%#02x '%s'", cell.Char, utils.Ternary(cell.Char >= 0x40, string(conv.E2A(cell.Char)), " "))

	// 👇 cell
	t.AppendRow(table.Row{
		char,
		l.boolean(cell.IsFldStart()),
		frow,
		fcol,
		crow,
		ccol,
		l.boolean(cell.Attrs.CharAttr),
		utils.Ternary(cell.Attrs.Autoskip, "SKIP", ""),
		utils.Ternary(cell.Attrs.Blink, "BLINK", ""),
		utils.Ternary(cell.Attrs.Color != 0x00, types.ColorFor(cell.Attrs.Color), ""),
		utils.Ternary(cell.Attrs.Hidden, "HIDDEN", ""),
		utils.Ternary(cell.Attrs.Highlight, "HILITE", ""),
		utils.Ternary(cell.Attrs.Intensify, "INTENSE", ""),
		utils.Ternary(cell.Attrs.MDT, "MDT", ""),
		utils.Ternary(cell.Attrs.Numeric, "NUM", ""),
		utils.Ternary(cell.Attrs.Protected, "PROT", ""),
		utils.Ternary(cell.Attrs.Reverse, "REV", ""),
		utils.Ternary(cell.Attrs.Underscore, "USCORE", ""),
		utils.Ternary(cell.Attrs.Outline != 0x00, types.OutlineFor(cell.Attrs.Outline), ""),
		utils.Ternary(cell.Attrs.LCID != 0x00, cell.Attrs.LCID.String(), ""),
	})
}

// ---------------------------------------------------------------------------
// 🟦 Trace pubsub activity
//
// 🔥 currently disabled
// ---------------------------------------------------------------------------

func (l *Logger) logTrace(topic Topic, handler interface{}) {
	if !strings.Contains(topic.String(), "tick") /* 🔥 suppressed ?? */ && false {
		pkg, nm := utils.GetFuncName(handler)
		fmt.Printf("🐞 topic %s -> func %s() in %s\n", topic, nm, pkg)
	}
}

// ---------------------------------------------------------------------------
// 🟦 Log the WCC
// ---------------------------------------------------------------------------

func (l *Logger) logWCC(wcc types.WCC) {
	t := l.newTable(text.FgHiBlue, "")
	defer t.Render()

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter},
		{Number: 4, Align: text.AlignCenter},
		{Number: 5, Align: text.AlignCenter},
	})

	t.AppendHeader(table.Row{
		"",
		"Alarm",
		"Reset",
		"ResetMDT",
		"Unlock",
	})

	t.AppendRow(table.Row{
		"WCC",
		l.boolean(wcc.Alarm),
		l.boolean(wcc.Reset),
		l.boolean(wcc.ResetMDT),
		l.boolean(wcc.Unlock),
	})
}

// ---------------------------------------------------------------------------
// 🟦 Helpers
// ---------------------------------------------------------------------------

func (l *Logger) boolean(flag bool) string {
	return utils.Ternary(flag, "\u2022", "")
}

func (l *Logger) newTable(color text.Color, title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	style := table.StyleBold
	style.Color = table.ColorOptions{
		Border:    text.Colors{color, text.Bold},
		Separator: text.Colors{color, text.Bold},
		Header:    text.Colors{color, text.Bold},
		Row:       text.Colors{text.Reset},
	}
	t.SetStyle(style)
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 1,
			Colors: text.Colors{color, text.Bold},
		},
	})
	if title != "" {
		t.Style().Title.Align = text.AlignCenter
		t.Style().Title.Colors = text.Colors{color, text.Bold}
		t.SetTitle(title)
	}
	return t
}

func (l *Logger) wrap(w int) text.Transformer {
	return func(val interface{}) string {
		return text.WrapText(fmt.Sprint(val), w)
	}
}

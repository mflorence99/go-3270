package debug

import (
	"fmt"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log cell contents (STRL+arrow)

func (l *Logger) logProbe(addr int) {
	t := l.newTable(text.FgHiRed, "")
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
	})
	t.AppendHeader(table.Row{
		"",
		"SF",
		"Row",
		"Col",
		"Row",
		"Col",
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
	cell, _ := l.buf.Peek(addr)
	crow, ccol := l.cfg.Addr2RC(addr)
	frow, fcol := l.cfg.Addr2RC(cell.FldAddr)
	char := fmt.Sprintf("%#02x '%s'", conv.A2E(cell.Char), utils.Ternary(cell.Char >= 0x20, string(cell.Char), " "))
	// 👇 cell
	t.AppendRow(table.Row{
		char,
		l.boolean(cell.FldStart),
		frow,
		fcol,
		crow,
		ccol,
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

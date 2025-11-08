package debug

import (
	"fmt"
	"go3270/emulator/conv"
	"go3270/emulator/utils"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log buffer contents

func (l *Logger) logBuffer() {

	cursor := fmt.Sprintf("%s%s", text.FgHiGreen.Sprint("\u25a0"), text.FgWhite.Sprint("\u200b"))

	protected := fmt.Sprintf("%s%s", text.FgHiCyan.Sprint("\u00b6"), text.FgWhite.Sprint("\u200b"))

	unprotected := fmt.Sprintf("%s%s", text.FgHiRed.Sprint("\u00b6"), text.FgWhite.Sprint("\u200b"))

	// 👇 define the table
	t := l.newTable(text.FgHiBlue, fmt.Sprintf("%s Buffer\nCURSOR: %s PROT: %s UNPROT: %s", l.buf.Mode(), cursor, protected, unprotected))
	defer t.Render()

	// 👇 header rows
	row1 := ""
	row2 := ""
	for ix := 10; ix <= l.cfg.Cols; ix += 10 {
		row1 += fmt.Sprintf("%10d", ix/10)
		row2 += "1234567890"
	}
	t.AppendHeader(table.Row{
		"",
		fmt.Sprintf("%s\n%s", row1, row2),
	})

	// 👇 where's the cursorAt?
	row, col := l.cfg.Addr2RC(l.st.Status.CursorAt)

	// 👇 data rows
	for iy := 1; iy <= l.cfg.Rows; iy++ {
		var b strings.Builder
		// 👇 data cols
		for ix := 1; ix <= l.cfg.Cols; ix++ {
			// 👇 show the cursor specially
			if iy == row && ix == col {
				b.WriteString(cursor)
			} else {
				// 👇 or the cell contents, best as we can
				cell, ok := l.buf.Peek(ix + ((iy - 1) * l.cfg.Cols) - 1)
				if cell != nil && ok {
					if cell.FldStart {
						b.WriteString(utils.Ternary(cell.Attrs.Protected, protected, unprotected))
					} else {
						b.WriteRune(utils.Ternary(cell.Char <= 0x40, ' ', conv.E2Rune(cell.Attrs.LCID, cell.Char)))
					}
				}
			}
		}
		t.AppendRow(table.Row{iy, b.String()})
	}
}

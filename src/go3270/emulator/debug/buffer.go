package debug

import (
	"fmt"
	"go3270/emulator/conv"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log buffer contents

func (l *Logger) logBuffer() {
	t := l.newTable(text.FgHiBlue, fmt.Sprintf("%s Buffer\nLEGEND: \u25a0 cursor -- \u00b6 protected -- \u00bb", l.buf.Mode()))
	defer t.Render()
	// 👇 header rows
	row1 := ""
	row2 := ""
	for ix := 10; ix <= l.cfg.Cols; ix += 10 {
		row1 += fmt.Sprintf("%10d", ix/10)
		row2 += "1234567890"
	}
	t.AppendHeader(table.Row{"", fmt.Sprintf("%s\n%s", row1, row2)})
	// 👇 where's the cursorAt?
	crow, ccol := l.cfg.Addr2RC(l.st.Status.CursorAt)
	// 👇 data rows
	for iy := 1; iy <= l.cfg.Rows; iy++ {
		row := ""
		for ix := 1; ix <= l.cfg.Cols; ix++ {
			// 👇 show the cursor specially
			if iy == crow && ix == ccol {
				row += "\u25a0"
			} else {
				// 👇 or the cell contents, best as we can
				cell, ok := l.buf.Peek(ix + ((iy - 1) * l.cfg.Cols) - 1)
				if cell != nil && ok {
					if cell.FldStart {
						row += utils.Ternary(cell.Attrs.Protected, "\u00b6", "\u00bb")
					} else {
						row += string(utils.Ternary(cell.Char <= 0x40, ' ', conv.E2A(cell.Char)))
					}
				}
			}
		}
		t.AppendRow(table.Row{iy, row})
	}
}

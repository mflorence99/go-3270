package debug

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
)

func LogBuffer(cfg pubsub.Config, buf *buffer.Buffer) {
	t := NewTable()
	defer t.Render()
	// 👇 header rows
	row1 := ""
	row2 := ""
	for ix := 10; ix <= cfg.Cols; ix += 10 {
		row1 += fmt.Sprintf("%10d", ix/10)
		row2 += "1234567890"
	}
	t.AppendHeader(table.Row{"", fmt.Sprintf("%s\n%s", row1, row2)})
	// 👇 data rows
	for ix := 1; ix <= cfg.Rows; ix++ {
		row := ""
		for addr := 0; addr < cfg.Cols; addr++ {
			cell, _ := buf.Peek(addr + ((ix - 1) * cfg.Cols))
			if cell.FldStart {
				row += utils.Ternary(cell.Attrs.Protected, "\u00b6", "\u00bb")
			} else {
				row += string(utils.Ternary(cell.Char <= ' ', ' ', cell.Char))
			}
		}
		t.AppendRow(table.Row{ix, row})
	}
}

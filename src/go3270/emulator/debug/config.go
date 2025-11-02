package debug

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log configuration

func (l *Logger) logConfig() {
	t := l.newTable(text.FgHiBlue, "Config")
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "", "", "BG", "", "Font", "Font", "Font", "Padding", "Padding"})
	t.AppendHeader(table.Row{"", "#Rows", "#Cols", "Color", "Mono", "Highlight", "Width", "Height", "Size", "Width", "Height"})
	t.AppendRow(table.Row{
		"CFG", l.cfg.Rows, l.cfg.Cols, l.cfg.BgColor, l.boolean(l.cfg.Monochrome), l.cfg.FontWidth, l.cfg.FontHeight, l.cfg.FontSize, l.cfg.PaddedWidth, l.cfg.PaddedHeight},
	)
}

func (l *Logger) logCLUT() {
	t := l.newTable(text.FgHiBlue, "CLUT")
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "Attr", "Color"})
	for k, v := range l.cfg.CLUT {
		t.AppendRow(table.Row{k, fmt.Sprintf("%#02x", byte(k)), v})
	}
}

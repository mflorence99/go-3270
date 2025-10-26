package debug

import (
	"fmt"
	"go3270/emulator/pubsub"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (l *Logger) logConfig(cfg pubsub.Config) {
	t := l.newTable(text.FgHiBlue)
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "", "", "BG", "", "Font", "Font", "Font", "Padding", "Padding"})
	t.AppendHeader(table.Row{"", "#Rows", "#Cols", "Color", "Mono", "Highlight", "Width", "Height", "Size", "Width", "Height"})
	t.AppendRow(table.Row{
		"CFG", cfg.Rows, cfg.Cols, cfg.BgColor, l.boolean(cfg.Monochrome), cfg.FontWidth, cfg.FontHeight, cfg.FontSize, cfg.PaddedWidth, cfg.PaddedHeight},
	)
}

func (l *Logger) logCLUT(cfg pubsub.Config) {
	t := l.newTable(text.FgHiBlue)
	defer t.Render()
	// 👇 table rows
	t.AppendHeader(table.Row{"", "Attr", "Color"})
	for k, v := range cfg.CLUT {
		t.AppendRow(table.Row{k, fmt.Sprintf("%#02x", byte(k)), v})
	}
}

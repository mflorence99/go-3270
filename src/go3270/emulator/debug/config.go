package debug

import (
	"go3270/emulator/pubsub"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (l *Logger) logConfig(cfg pubsub.Config) {
	t := NewTable()
	defer t.Render()
	t.AppendHeader(table.Row{"", "", "Color", "Color", "Color", "Font", "Font", "Font", "Padding", "Padding"})
	t.AppendHeader(table.Row{"#Rows", "#Cols", "BG", "Normal", "Highlight", "Width", "Height", "Size", "Width", "Height"})
	t.AppendRow(table.Row{
		cfg.Rows, cfg.Cols, cfg.BgColor, cfg.Color[0], cfg.Color[1], cfg.FontWidth, cfg.FontHeight, cfg.FontSize, cfg.PaddedWidth, cfg.PaddedHeight},
	)
}

func (l *Logger) logCLUT(cfg pubsub.Config) {
	t := NewTable()
	defer t.Render()
	t.AppendHeader(table.Row{"", "Color", "Color"})
	t.AppendHeader(table.Row{"", "Normal", "Highlight"})
	for k, v := range cfg.CLUT {
		t.AppendRow(table.Row{k, v[0], v[1]})
	}
}

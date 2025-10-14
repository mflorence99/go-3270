package screen

import (
	"go3270/emulator/pubsub"
	"math"
)

type Screen struct {
	Cells []Box
}

func NewScreen(cfg pubsub.Config) *Screen {
	s := new(Screen)
	s.Cells = make([]Box, cfg.Cols*cfg.Rows)
	for ix := range s.Cells {
		w := math.Round(cfg.FontWidth * cfg.PaddedWidth)
		h := math.Round(cfg.FontHeight * cfg.PaddedHeight)
		col := ix % cfg.Cols
		row := int(ix / cfg.Cols)
		x := math.Round(float64(col) * w)
		y := -math.Round(float64(row) * h)
		// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
		baseline := y + h - (cfg.FontSize / 2)
		s.Cells[ix] = Box{x, y, w, h, baseline}
	}
	return s
}

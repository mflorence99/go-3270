package pubsub

import (
	"fmt"
	"math"
)

type Box struct {
	X        float64
	Y        float64
	W        float64
	H        float64
	Baseline float64
}

func NewBox(row, col int, cfg Config) Box {
	w := math.Round(cfg.FontWidth * cfg.PaddedWidth)
	h := math.Round(cfg.FontHeight * cfg.PaddedHeight)
	x := math.Round(float64(col) * w)
	y := math.Round(float64(row) * h)
	// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
	baseline := y + h - (cfg.FontSize / 2)
	return Box{x, y, w, h, baseline}
}

func (b Box) String() string {
	return fmt.Sprintf("xyb[%f %f %f] wh[%f %f] ", b.X, b.Y, b.Baseline, b.W, b.H)
}

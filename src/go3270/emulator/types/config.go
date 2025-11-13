package types

import (
	"image"

	"golang.org/x/image/font"
)

// 🟧 Go3270 configuration parameters

type Config struct {
	BgColor      string
	BoldFace     *font.Face
	CLUT         map[Color]string
	Cols         int
	FontHeight   float64
	FontSize     float64
	FontWidth    float64
	Monochrome   bool
	NormalFace   *font.Face
	PaddedHeight float64
	PaddedWidth  float64
	RGBA         *image.RGBA
	Rows         int
	Screenshot   string
}

// 🟦 Public functions

func (c Config) Addr2RC(addr int) (int, int) {
	row := int(addr/c.Cols) + 1
	col := (addr % c.Cols) + 1
	return row, col
}

func (c Config) RC2Addr(row, col int) int {
	return (row-1)*c.Cols + c.Cols - 1
}

func (c Config) ColorOf(a *Attrs) string {
	var ix Color
	if c.Monochrome {
		ix = 0xf4
	} else if a.Color == 0x00 {
		switch {
		case !a.Protected && (a.Highlight || a.Hidden):
			ix = 0xF2
		case !a.Protected && !a.Highlight:
			ix = 0xF4
		case a.Protected && (a.Highlight || a.Hidden):
			ix = 0xF7
		case a.Protected && !a.Highlight:
			ix = 0xF1
		default:
			ix = 0xF4
		}
	} else {
		ix = a.Color
	}
	return c.CLUT[ix]
}

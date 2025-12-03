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
	Cols         uint
	FontHeight   float64
	FontSize     float64
	FontWidth    float64
	Monochrome   bool
	NormalFace   *font.Face
	PaddedHeight float64
	PaddedWidth  float64
	RGBA         *image.RGBA
	Rows         uint
	SuppressLogs bool
	Testpage     string
}

// 🟦 Public functions

func (c *Config) Addr2RC(addr uint) (uint, uint) {
	row := uint(addr/c.Cols) + 1
	col := (addr % c.Cols) + 1
	return row, col
}

func (c *Config) RC2Addr(row, col uint) uint {
	return (row-1)*c.Cols + col - 1
}

func (c *Config) ColorOf(a *Attrs) string {
	var ix Color
	if c.Monochrome {
		ix = GREEN
	} else if a.Color == 0x00 {
		switch {
		case !a.Protected && (a.Highlight || a.Hidden):
			ix = RED
		case !a.Protected && !a.Highlight:
			ix = GREEN
		case a.Protected && (a.Highlight || a.Hidden):
			ix = FOREGROUND
		case a.Protected && !a.Highlight:
			ix = BLUE
		default:
			ix = 0xf4
		}
	} else {
		ix = a.Color
	}
	return c.CLUT[ix]
}

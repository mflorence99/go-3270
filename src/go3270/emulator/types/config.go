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

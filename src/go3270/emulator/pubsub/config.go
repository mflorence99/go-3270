package pubsub

import (
	"go3270/emulator/consts"
	"image"

	"golang.org/x/image/font"
)

type Config struct {
	BgColor      string
	BoldFace     *font.Face
	CLUT         map[consts.Color][2]string
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

func (c Config) Addr2RC(addr int) (int, int) {
	row := int(addr/c.Cols) + 1
	col := (addr % c.Cols) + 1
	return row, col
}

func (c Config) RC2Addr(row, col int) int {
	return (row-1)*c.Cols + c.Cols - 1
}

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
	Color        [2]string
	Cols         int
	FontHeight   float64
	FontSize     float64
	FontWidth    float64
	NormalFace   *font.Face
	PaddedHeight float64
	PaddedWidth  float64
	RGBA         *image.RGBA
	Rows         int
}

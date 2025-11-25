package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheImageFor(t *testing.T) {
	emu := MockEmulator().Init()

	g := Glyph{
		Char:      'A',
		Color:     "#123456",
		Highlight: true,
		LCID:      0xf1,
		Outline: Outline{
			Bottom: true,
			Right:  true,
			Top:    true,
			Left:   true,
		},
		Reverse:    true,
		Underscore: true,
	}

	img1 := emu.GC.ImageFor(g, NewBox(1, 2, emu.Cfg))
	img2 := emu.GC.ImageFor(g, NewBox(3, 4, emu.Cfg))
	assert.Equal(t, img1, img2, "cache returns identical image")

	g.Outline.Left = false
	img3 := emu.GC.ImageFor(g, NewBox(5, 6, emu.Cfg))
	assert.NotEqual(t, img1, img3, "different glyphs create different images")
}

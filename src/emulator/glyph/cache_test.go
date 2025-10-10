package glyph_test

import (
	"emulator/glyph"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetSet(t *testing.T) {
	cache := new(glyph.Cache)
	glyph := glyph.Glyph{
		Char:       0x00,
		Color:      "#999999",
		Highlight:  true,
		Reverse:    true,
		Underscore: true,
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 100, 100))
	img, ok := cache.Get(glyph)
	assert.True(t, img == nil)
	assert.False(t, ok)
	cache.Set(glyph, rgba)
	img, ok = cache.Get(glyph)
	assert.True(t, img != nil)
	assert.True(t, ok)
}

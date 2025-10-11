package glyph_test

import (
	"emulator/glyph"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	g := glyph.NewGlyph()
	g.Char = 0x00
	g.Color = "#999999"
	g.Highlight = true
	g.Reverse = true
	g.Underscore = true
	str := g.String()
	assert.True(t, str == "GLYPH=[ 0x00 #999999 HILITE REV USCORE ]")
}

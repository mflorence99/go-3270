package glyph_test

import (
	"emulator/glyph"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	glyph := glyph.Glyph{
		Char:       0x00,
		Color:      "#999999",
		Highlight:  true,
		Reverse:    true,
		Underscore: true,
	}
	str := glyph.String()
	assert.True(t, str == "GLYPH=[ 0x00 #999999 HILITE REV USCORE ]")
}

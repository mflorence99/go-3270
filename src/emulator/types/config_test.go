package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddr2RC(t *testing.T) {
	c := &Config{Cols: 80, Rows: 24}
	row, col := c.Addr2RC(100)
	assert.Equal(t, row, uint(2), "row calculated from addr")
	assert.Equal(t, col, uint(21), "col calculated from addr")
}

func TestRC2Addr(t *testing.T) {
	c := &Config{Cols: 80, Rows: 24}
	addr := c.RC2Addr(uint(2), uint(21))
	assert.Equal(t, uint(100), addr, "addr calculated from row/col")
}

func TestColorOf(t *testing.T) {
	c := &Config{CLUT: map[Color]string{
		0xf1: "WHITE",
		0xf2: "RED",
		0xf4: "GREEN",
		0xf7: "BLUE",
	}}

	a := &Attrs{Protected: true}
	assert.Equal(t, "WHITE", c.ColorOf(a), "color from protected")

	a = &Attrs{Hidden: true}
	assert.Equal(t, "RED", c.ColorOf(a), "color from hidden")

	a = &Attrs{}
	assert.Equal(t, "GREEN", c.ColorOf(a), "default color")

	a = &Attrs{Protected: true, Highlight: true}
	assert.Equal(t, "BLUE", c.ColorOf(a), "color from protected+highlight")

}

func TestColorOfMono(t *testing.T) {
	c := &Config{Monochrome: true, CLUT: map[Color]string{0xf4: "GREEN"}}
	a := &Attrs{}
	assert.Equal(t, "GREEN", c.ColorOf(a), "monochrome display is 'green'")
}

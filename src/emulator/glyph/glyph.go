package glyph

import (
	"fmt"
)

type Glyph struct {
	Char       byte
	Color      string
	Highlight  bool
	Reverse    bool
	Underscore bool
}

func NewGlyph() *Glyph {
	g := new(Glyph)
	return g
}

func (g *Glyph) String() string {
	str := fmt.Sprintf("GLYPH=[ 0x%02x %s ", g.Char, g.Color)
	if g.Highlight {
		str += "HILITE "
	}
	if g.Reverse {
		str += "REV "
	}
	if g.Underscore {
		str += "USCORE "
	}
	str += "]"
	return str
}

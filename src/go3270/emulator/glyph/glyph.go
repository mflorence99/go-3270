package glyph

import (
	"fmt"
)

// 🟧 A glyph, as stored in the glyph cache

type Glyph struct {
	Char       byte
	Color      string
	Highlight  bool
	Reverse    bool
	Underscore bool
}

// 🟦 Constructor

func NewGlyph() *Glyph {
	g := new(Glyph)
	return g
}

// 🟦 Stringer implementation

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

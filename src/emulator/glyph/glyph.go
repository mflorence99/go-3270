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

func (glyph *Glyph) String() string {
	str := fmt.Sprintf("GLYPH=[ 0x%02x %s ", glyph.Char, glyph.Color)
	if glyph.Highlight {
		str += "HILITE "
	}
	if glyph.Reverse {
		str += "REV "
	}
	if glyph.Underscore {
		str += "USCORE "
	}
	str += "]"
	return str
}

package screen

import "go3270/emulator/types"

// 🟧 A glyph, as stored in the glyph cache

type Glyph struct {
	Char       byte
	Color      string
	Highlight  bool
	Reverse    bool
	Underscore bool
	Outline    Outline
	LCID       types.LCID
}

type Outline struct {
	Bottom bool
	Right  bool
	Top    bool
	Left   bool
}

// 🟦 Constructor

func NewGlyph() *Glyph {
	g := new(Glyph)
	return g
}

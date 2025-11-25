package core

import "emulator/types"

// 🟧 A glyph, as stored in the glyph cache

type Glyph struct {
	Char       byte
	Color      string
	Highlight  bool
	LCID       types.LCID
	Outline    Outline
	Reverse    bool
	Underscore bool
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

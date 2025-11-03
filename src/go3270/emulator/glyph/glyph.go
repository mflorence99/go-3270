package glyph

// 🟧 A glyph, as stored in the glyph cache

type Glyph struct {
	Char       byte
	Color      string
	Highlight  bool
	Reverse    bool
	Underscore bool
	Outline    struct {
		Bottom bool
		Right  bool
		Top    bool
		Left   bool
	}
}

// 🟦 Constructor

func NewGlyph() *Glyph {
	g := new(Glyph)
	return g
}

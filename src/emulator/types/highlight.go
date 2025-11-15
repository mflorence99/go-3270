package types

// 🟧 3270 highlight extended attribute

type Highlight byte

// 🟦 Lookup tables

const (
	DFLT_HILITE Highlight = 0x00
	NO_HILITE   Highlight = 0xF0
	BLINK       Highlight = 0xF1
	REVERSE     Highlight = 0xF2
	UNDERSCORE  Highlight = 0xF4
	INTENSIFY   Highlight = 0xF8
)

var highlights = map[Highlight]string{
	0x00: "DFLT_HILITE",
	0xF0: "NO_HILITE",
	0xF1: "BLINK",
	0xF2: "REVERSE",
	0xF4: "UNDERSCORE",
	0xF8: "INTENSIFY",
}

// 🟦 Stringer implementation

func HighlightFor(h Highlight) string {
	return highlights[h]
}

func (h Highlight) String() string {
	return HighlightFor(h)
}

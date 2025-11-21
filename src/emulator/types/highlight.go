package types

// 🟧 3270 highlight extended attribute

type Highlight byte

// 🟦 Lookup tables

const (
	DFLT_HILITE Highlight = 0x00
	NO_HILITE   Highlight = 0xf0
	BLINK       Highlight = 0xf1
	REVERSE     Highlight = 0xf2
	UNDERSCORE  Highlight = 0xf4
	INTENSIFY   Highlight = 0xf8
)

var highlights = map[Highlight]string{
	0x00: "DFLT_HILITE",
	0xf0: "NO_HILITE",
	0xf1: "BLINK",
	0xf2: "REVERSE",
	0xf4: "UNDERSCORE",
	0xf8: "INTENSIFY",
}

// 🟦 Stringer implementation

func HighlightFor(h Highlight) string {
	return highlights[h]
}

func (h Highlight) String() string {
	return HighlightFor(h)
}

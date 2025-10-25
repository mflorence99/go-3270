package consts

type Highlight byte

const (
	NOHILITE   Highlight = 0xF0
	BLINK      Highlight = 0xF1
	REVERSE    Highlight = 0xF2
	UNDERSCORE Highlight = 0xF4
	INTENSIFY  Highlight = 0xF8
)

var highlights = map[Highlight]string{
	0xF0: "NOHILITE",
	0xF1: "BLINK",
	0xF2: "REVERSE",
	0xF4: "UNDERSCORE",
	0xF8: "INTENSITY",
}

func HighlightFor(h Highlight) string {
	return highlights[h]
}

func (h Highlight) String() string {
	return HighlightFor(h)
}

package consts

type Highlight byte

const (
	BLINK      Highlight = 0xF1
	REVERSE    Highlight = 0xF2
	UNDERSCORE Highlight = 0xF4
)

var highlights = map[Highlight]string{
	0xF1: "BLINK",
	0xF2: "REVERSE",
	0xF4: "UNDERSCORE",
}

func HighlightFor(h Highlight) string {
	return highlights[h]
}

func (h Highlight) String() string {
	return HighlightFor(h)
}

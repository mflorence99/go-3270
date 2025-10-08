package consts

var (
	BLINK      byte = 0xF1
	REVERSE    byte = 0xF2
	UNDERSCORE byte = 0xF4
)

var highlights = map[byte]string{
	0xF1: "BLINK",
	0xF2: "REVERSE",
	0xF4: "UNDERSCORE",
}

func HighlightFor(highlight byte) string {
	return highlights[highlight]
}

package consts

// 🟧 3270 colors

type Color byte

// 🟦 Lookup tables

const (
	BACKGROUND     Color = 0xF0
	BLUE           Color = 0xF1
	RED            Color = 0xF2
	PINK           Color = 0xF3
	GREEN          Color = 0xF4
	TURQUOISE      Color = 0xF5
	YELLOW         Color = 0xF6
	FOREGROUND     Color = 0xF7
	BLACK          Color = 0xF8
	DEEP_BLUE      Color = 0xF9
	ORANGE         Color = 0xFA
	PURPLE         Color = 0xFB
	PALE_GREEN     Color = 0xFC
	PALE_TURQUOISE Color = 0xFD
	GREY           Color = 0xFE
	WHITE          Color = 0xFF
)

var colors = map[Color]string{
	0xF0: "BACKGROUND",
	0xF1: "BLUE",
	0xF2: "RED",
	0xF3: "PINK",
	0xF4: "GREEN",
	0xF5: "TURQUOISE",
	0xF6: "YELLOW",
	0xF7: "FOREGROUND",
	0xF8: "BLACK",
	0xF9: "DEEP_BLUE",
	0xFA: "ORANGE",
	0xFB: "PURPLE",
	0xFC: "PALE_GREEN",
	0xFD: "PALE_TURQUOISE",
	0xFE: "GREY",
	0xFF: "WHITE",
}

// 🟦 Stringer implementation

func ColorFor(c Color) string {
	return colors[c]
}

func (c Color) String() string {
	return ColorFor(c)
}

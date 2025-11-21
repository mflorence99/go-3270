package types

// 🟧 3270 colors

type Color byte

// 🟦 Lookup tables

const (
	BACKGROUND     Color = 0xf0
	BLUE           Color = 0xf1
	RED            Color = 0xf2
	PINK           Color = 0xf3
	GREEN          Color = 0xf4
	TURQUOISE      Color = 0xf5
	YELLOW         Color = 0xf6
	FOREGROUND     Color = 0xf7
	BLACK          Color = 0xf8
	DEEP_BLUE      Color = 0xf9
	ORANGE         Color = 0xfA
	PURPLE         Color = 0xfB
	PALE_GREEN     Color = 0xfC
	PALE_TURQUOISE Color = 0xfD
	GREY           Color = 0xfE
	WHITE          Color = 0xfF
)

var colors = map[Color]string{
	0xf0: "BACKGROUND",
	0xf1: "BLUE",
	0xf2: "RED",
	0xf3: "PINK",
	0xf4: "GREEN",
	0xf5: "TURQUOISE",
	0xf6: "YELLOW",
	0xf7: "FOREGROUND",
	0xf8: "BLACK",
	0xf9: "DEEP_BLUE",
	0xfA: "ORANGE",
	0xfB: "PURPLE",
	0xfC: "PALE_GREEN",
	0xfD: "PALE_TURQUOISE",
	0xfE: "GREY",
	0xfF: "WHITE",
}

// 🟦 Stringer implementation

func ColorFor(c Color) string {
	return colors[c]
}

func (c Color) String() string {
	return ColorFor(c)
}

package consts

// 👁️ https://bitsavers.trailing-edge.com/pdf/ibm/3278/GA33-3056-0_3270_Information_Display_System_Color_and_Programmed_Symbols_3278_3279_3287_Sep1979.pdf?utm_source=chatgpt.com

type Color byte

const (
	BLACK     Color = 0xF0
	BLUE      Color = 0xF1
	RED       Color = 0xF2
	PINK      Color = 0xF3
	GREEN     Color = 0xF4
	TURQUOISE Color = 0xF5
	YELLOW    Color = 0xF6
	WHITE     Color = 0xF7
)

var colors = map[Color]string{
	0xF0: "BLACK",
	0xF1: "BLUE",
	0xF2: "RED",
	0xF3: "PINK",
	0xF4: "GREEN",
	0xF5: "TURQUOISE",
	0xF6: "YELLOW",
	0xF7: "WHITE",
}

func ColorFor(c Color) string {
	return colors[c]
}

func (c Color) String() string {
	return ColorFor(c)
}

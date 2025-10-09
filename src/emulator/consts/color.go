package consts

// 👁️ https://bitsavers.trailing-edge.com/pdf/ibm/3278/GA33-3056-0_3270_Information_Display_System_Color_and_Programmed_Symbols_3278_3279_3287_Sep1979.pdf?utm_source=chatgpt.com

type Color [2]string

var (
	BLACK     byte = 0xF0
	BLUE      byte = 0xF1
	RED       byte = 0xF2
	PINK      byte = 0xF3
	GREEN     byte = 0xF4
	TURQUOISE byte = 0xF5
	YELLOW    byte = 0xF6
	WHITE     byte = 0xF7
)

var colors = map[byte]string{
	0xF0: "BLACK",
	0xF1: "BLUE",
	0xF2: "RED",
	0xF3: "PINK",
	0xF4: "GREEN",
	0xF5: "TURQUOISE",
	0xF6: "YELLOW",
	0xF7: "WHITE",
}

var CLUT = map[byte]Color{
	0xF0: {"#111138", "#505050"},
	0xF1: {"#0078FF", "#3366CC"},
	0xF2: {"#D40000", "#E06666"},
	0xF3: {"#FF69B4", "#FFB3DA"},
	0xF4: {"#00AA00", "#88DD88"},
	0xF5: {"#00C8AA", "#99E8DD"},
	0xF6: {"#FF8000", "#FFB266"},
	0xF7: {"#888888", "#FFFFFF"},
}

func ColorFor(color byte) string {
	return colors[color]
}

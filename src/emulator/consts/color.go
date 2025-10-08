package consts

type Color [2]string

// 👁️ https://bitsavers.trailing-edge.com/pdf/ibm/3278/GA33-3056-0_3270_Information_Display_System_Color_and_Programmed_Symbols_3278_3279_3287_Sep1979.pdf?utm_source=chatgpt.com
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

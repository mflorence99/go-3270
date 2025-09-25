package main

// 👁️ https://bitsavers.trailing-edge.com/pdf/ibm/3278/GA33-3056-0_3270_Information_Display_System_Color_and_Programmed_Symbols_3278_3279_3287_Sep1979.pdf?utm_source=chatgpt.com
var CLUT = map[int][]string{
	0xf0: {"#111138", "#505050"},
	0xf1: {"#0078FF", "#3366CC"},
	0xf2: {"#D40000", "#E06666"},
	0xf3: {"#FF69B4", "#FFB3DA"},
	0xf4: {"#00AA00", "#88DD88"},
	0xf5: {"#00C8AA", "#99E8DD"},
	0xf6: {"#FF8000", "#FFB266"},
	0xf7: {"#FFFFFF", "#B8B8B8"},
}

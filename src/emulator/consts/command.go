package consts

var (
	EAU byte = 0x6F
	EW  byte = 0xF5
	EWA byte = 0x7E
	RB  byte = 0xF2
	RM  byte = 0xF6
	RMA byte = 0x6E
	W   byte = 0xF1
	WSF byte = 0xF3
)

var commands = map[byte]string{0x6E: "RMA",
	0x6F: "EAU",
	0x7E: "EWA",
	0xF1: "W",
	0xF2: "RB",
	0xF3: "WSF",
	0xF5: "EW",
	0xF6: "RM",
}

func CommandFor(command byte) string {
	return commands[command]
}

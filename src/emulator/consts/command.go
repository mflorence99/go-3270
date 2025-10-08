package consts

var EAU byte = 0x6F
var EW byte = 0xF5
var EWA byte = 0x7E
var RB byte = 0xF2
var RM byte = 0xF6
var RMA byte = 0x6E
var W byte = 0xF1
var WSF byte = 0xF3

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

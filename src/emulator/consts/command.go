package consts

type Command byte

const (
	EAU Command = 0x6F
	EW  Command = 0xF5
	EWA Command = 0x7E
	Q   Command = 0x02
	QL  Command = 0x03
	RB  Command = 0xF2
	RM  Command = 0xF6
	RMA Command = 0x6E
	W   Command = 0xF1
	WSF Command = 0xF3
)

var commands = map[Command]string{
	0x02: "Q",
	0x03: "QL",
	0x6E: "RMA",
	0x6F: "EAU",
	0x7E: "EWA",
	0xF1: "W",
	0xF2: "RB",
	0xF3: "WSF",
	0xF5: "EW",
	0xF6: "RM",
}

func CommandFor(command Command) string {
	return commands[command]
}

func (command Command) String() string {
	return CommandFor(command)
}

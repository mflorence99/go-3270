package types

// 🟧 3270 Commands

type Command byte

// 🟦 Lookup tables

const (
	EAU Command = 0x6F
	EW  Command = 0xf5
	EWA Command = 0x7E
	Q   Command = 0x02
	QL  Command = 0x03
	RB  Command = 0xf2
	RM  Command = 0xf6
	RMA Command = 0x6E
	W   Command = 0xf1
	WSF Command = 0xf3
)

var commands = map[Command]string{
	0x02: "Q",
	0x03: "QL",
	0x6E: "RMA",
	0x6F: "EAU",
	0x7E: "EWA",
	0xf1: "W",
	0xf2: "RB",
	0xf3: "WSF",
	0xf5: "EW",
	0xf6: "RM",
}

// 🟦 Stringer implementation

func CommandFor(c Command) string {
	return commands[c]
}

func (c Command) String() string {
	return CommandFor(c)
}

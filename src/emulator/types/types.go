package types

// 🟦 AID

var AID = map[uint8]string{
	0x88: "DEFAULT",
	0x6D: "CLEAR",
	0x7D: "ENTER",
	0x6C: "PA1",
	0x6E: "PA2",
	0x6B: "PA3",
	0xF1: "PF1",
	0xF2: "PF2",
	0xF3: "PF3",
	0xF4: "PF4",
	0xF5: "PF5",
	0xF6: "PF6",
	0xF7: "PF7",
	0xF8: "PF8",
	0xF9: "PF9",
	0x7A: "PF10",
	0x7B: "PF11",
	0x7C: "PF12",
	0xC1: "PF13",
	0xC2: "PF14",
	0xC3: "PF15",
	0xC4: "PF16",
	0xC5: "PF17",
	0xC6: "PF18",
	0xC7: "PF19",
	0xC8: "PF20",
	0xC9: "PF21",
	0x4A: "PF22",
	0x4B: "PF23",
	0x4C: "PF24",
}

var AIDLookup = make(map[string]uint8)

// 🟦 CLUT

// 👁️ https://bitsavers.trailing-edge.com/pdf/ibm/3278/GA33-3056-0_3270_Information_Display_System_Color_and_Programmed_Symbols_3278_3279_3287_Sep1979.pdf?utm_source=chatgpt.com
var CLUT = map[uint8][]string{
	0xf0: {"#111138", "#505050"},
	0xf1: {"#0078FF", "#3366CC"},
	0xf2: {"#D40000", "#E06666"},
	0xf3: {"#FF69B4", "#FFB3DA"},
	0xf4: {"#00AA00", "#88DD88"},
	0xf5: {"#00C8AA", "#99E8DD"},
	0xf6: {"#FF8000", "#FFB266"},
	0xf7: {"#FFFFFF", "#B8B8B8"},
}

// 🟦 Command

var Command = map[uint8]string{
	0x6F: "EAU",
	0x7E: "EWA",
	0xF1: "W ",
	0xF3: "WSF",
	0xF5: "EW",
}

var CommandLookup = make(map[string]uint8)

// 🟦 Op

var Op = map[uint8]string{
	0x02: "Q",
	0x03: "QL",
	0x6E: "RMA",
	0xF2: "RB",
	0xF6: "RM",
	0xFF: "UNKNOWN",
}

var OpLookup = make(map[string]uint8)

// 🟦 Order

var Order = map[uint8]string{
	0x05: "PT",
	0x08: "GE",
	0x11: "SBA",
	0x12: "EUA",
	0x13: "IC",
	0x1D: "SF",
	0x28: "SA",
	0x29: "SFE",
	0x2C: "MF",
	0x3C: "RA",
}

var OrderLookup = make(map[string]uint8)

// 🟧 Global Initialization (runs before main)

func init() {
	origs := []map[uint8]string{AID, Command, Op, Order}
	reverseds := []map[string]uint8{AIDLookup, CommandLookup, OpLookup, OrderLookup}
	for ix := 0; ix < len(origs); ix++ {
		orig := origs[ix]
		reversed := reverseds[ix]
		for k, v := range orig {
			reversed[v] = k
		}
	}
}

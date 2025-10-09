package device

// 🟦 AID

var AID = map[byte]string{
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

var AIDLookup = make(map[string]byte)

// 🟦 Command

var Command = map[byte]string{
	0x6E: "RMA",
	0x6F: "EAU",
	0x7E: "EWA",
	0xF1: "W",
	0xF2: "RB",
	0xF3: "WSF",
	0xF5: "EW",
	0xF6: "RM",
}

var CommandLookup = make(map[string]byte)

// 🟦 Highlight (attribute)

var Highlight = map[byte]string{
	0xF1: "BLINK",
	0xF2: "REVERSE",
	0xF4: "UNDERSCORE",
}

var HighlightLookup = make(map[string]byte)

// 🟦 LT (delineates outbound stream)

var FrameLT = []byte{0xFF, 0xEF}

// 🟦 Op

var Op = map[byte]string{
	0x02: "Q",
	0x03: "QL",
	0x6E: "RMA",
	0xF2: "RB",
	0xF6: "RM",
	0xFF: "UNKNOWN",
}

var OpLookup = make(map[string]byte)

// 🟦 Order

var Order = map[byte]string{
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

var OrderLookup = make(map[string]byte)

// 🟦 TypeCode (attribute)

var TypeCode = map[byte]string{
	0xC0: "BASIC",
	0x41: "HIGHLIGHT",
	0x42: "COLOR",
}

var TypeCodeLookup = make(map[string]byte)

// 🟧 Global Initialization (runs before main)

func init() {
	origs := []map[byte]string{AID, Command, Highlight, Op, Order, TypeCode}
	reverseds := []map[string]byte{AIDLookup, CommandLookup, HighlightLookup, OpLookup, OrderLookup, TypeCodeLookup}
	for ix := 0; ix < len(origs); ix++ {
		orig := origs[ix]
		reversed := reverseds[ix]
		for k, v := range orig {
			reversed[v] = k
		}
	}
}

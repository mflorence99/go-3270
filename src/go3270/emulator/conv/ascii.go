package conv

// 🟧  ASCII -> EBCDIC cconversion

// 🟦 Lookup tables

var ASCII = [256]byte{}

func init() {
	// 👇 the EBCDIC table starts at 0x40 because of the difficulty of eyeballing 64 0x00's
	for ix := 0; ix < len(EBCDIC); ix++ {
		ASCII[EBCDIC[ix]] = byte(ix + 0x40)
	}
}

// 🟦 Public functions

func A2E(a byte) byte {
	e := byte(0x00)
	if a == ' ' {
		e = 0x40
	} else {
		e = ASCII[a]
	}
	return e
}

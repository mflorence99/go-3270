package utils

// 🟦 Conversion between EBCDIC and ASCII

// 👁️ https://stackoverflow.com/questions/25367120/example-ebcdic-file-for-java-program-to-convert-ebcdic-to-ascii

// 🔥 by skipping the first 64 entries and starting on line 64, it's easy to read the EBCDIC character as the line # and the constant as its ASCII equivalent

var EBCDIC = []byte{
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	// start on line 64 to make reconciliation easier
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'¢',
	'.',
	'<',
	'(',
	'+',
	'|',
	'&',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'!',
	'$',
	'*',
	')',
	';',
	'¬',
	'-',
	'/',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'|',
	',',
	'%',
	'_',
	'>',
	'?',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'`',
	':',
	'#',
	'@',
	'\'',
	'=',
	'"',
	' ',
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
	'g',
	'h',
	'i',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'j',
	'k',
	'l',
	'm',
	'n',
	'o',
	'p',
	'q',
	'r',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	's',
	't',
	'u',
	'v',
	'w',
	'x',
	'y',
	'z',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'`',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'{',
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'G',
	'H',
	'I',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'}',
	'J',
	'K',
	'L',
	'M',
	'N',
	'O',
	'P',
	'Q',
	'R',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'\\',
	' ',
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	' ',
	' ',
	' ',
	' ',
	' ',
	' ',
}

var ASCII = [256]byte{}

func A2E(a []byte) []byte {
	e := make([]byte, len(a))
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == ' ' {
			e[ix] = 0x40
		} else {
			e[ix] = ASCII[a[ix]]
		}
	}
	return e
}

func E2A(e []byte) []byte {
	a := make([]byte, len(e))
	for ix := 0; ix < len(e); ix++ {
		if e[ix] >= 64 {
			a[ix] = EBCDIC[e[ix]-64]
		} else {
			a[ix] = ' '
		}
	}
	return a
}

// 🟦 3270 buffer address conversions

var Six2E = []byte{
	0x40, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7,
	0xC8, 0xC9, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
	0x50, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7,
	0xD8, 0xD9, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
	0x60, 0x61, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7,
	0xE8, 0xE9, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
	0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7,
	0xF8, 0xF9, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
}

func AddrFromBytes(u8s []byte) int {
	addr := int(u8s[0])
	addr &= 0b00111111
	addr = addr << 6
	addr += int(u8s[1] & 0b00111111)
	return addr
}

func AddrToBytes(addr int) []byte {
	u8s := make([]byte, 2)
	u8s[0] = Six2E[(addr>>6)&0b00111111]
	u8s[1] = Six2E[addr&0b00111111]
	return u8s
}

// 🟧 Global Initialization (runs before main)

func init() {
	// 👇 the EBCDIC table starts at 0x40 because of the difficulty of eyeballing 64 0x00's
	for ix := 0; ix < len(EBCDIC); ix++ {
		ASCII[EBCDIC[ix]] = byte(ix + 0x40)
	}
}

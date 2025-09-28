package utils

// 🟦 Conversion between EBCDIC and ASCII

// 👁️ https://stackoverflow.com/questions/25367120/example-ebcdic-file-for-java-program-to-convert-ebcdic-to-ascii

var EBCDIC = []uint8{
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

var ASCII = [256]uint8{}

func A2E(a []uint8) []uint8 {
	e := make([]uint8, len(a))
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == ' ' {
			e[ix] = 0x40
		} else {
			e[ix] = ASCII[a[ix]]
		}
	}
	return e
}

func E2A(e []uint8) []uint8 {
	a := make([]uint8, len(e))
	for ix := 0; ix < len(e); ix++ {
		if e[ix] >= 64 {
			a[ix] = EBCDIC[e[ix]-64]
		} else {
			a[ix] = ' '
		}
	}
	return a
}

// 🟧 Global Initialization (runs before main)

func init() {
	// 👇 the EBCDIC table starts at 0x40 because of the difficulty of eyeballing 64 0x00's
	for ix := 0; ix < len(EBCDIC); ix++ {
		ASCII[EBCDIC[ix]] = uint8(ix + 0x40)
	}
}

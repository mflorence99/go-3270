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
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'¢',
	'.',
	'<',
	'(',
	'+',
	'|',
	'&',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'!',
	'$',
	'*',
	')',
	';',
	'¬',
	'-',
	'/',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'|',
	',',
	'%',
	'_',
	'>',
	'?',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'`',
	':',
	'#',
	'@',
	'\'',
	'=',
	'"',
	0xA0,
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
	'g',
	'h',
	'i',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'j',
	'k',
	'l',
	'm',
	'n',
	'o',
	'p',
	'q',
	'r',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	's',
	't',
	'u',
	'v',
	'w',
	'x',
	'y',
	'z',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'`',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
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
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
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
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	'\\',
	0xA0,
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
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
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
	0xA0,
}

var ASCII = [256]uint8{}

// 🟧 Global Initialization (runs before main)

func init() {
	// 👇 the EBCDIC table starts at 0x40 because of the difficulty of eyeballing 64 0x00's
	for ix := 0; ix < len(EBCDIC); ix++ {
		ASCII[EBCDIC[ix]] = uint8(ix + 0x40)
	}
}

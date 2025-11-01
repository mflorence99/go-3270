package conv

// 🟧 EBCDIC -> ASCII conversion

// 👁️ https://stackoverflow.com/questions/25367120/example-ebcdic-file-for-java-program-to-convert-ebcdic-to-ascii

// 🔥 by skipping the first 64 entries and starting on line 64, it's easy to read the EBCDIC character as the line # and the constant as its ASCII equivalent

// 🟦 Lookup tables

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

// 🟦 Public functions

func E2A(e byte) byte {
	a := byte(0x00)
	if e >= 64 {
		a = EBCDIC[e-64]
	} else {
		a = ' '
	}
	return a
}

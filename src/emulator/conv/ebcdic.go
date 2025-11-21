package conv

// 🟧 EBCDIC -> ASCII conversion

// 🔥 by skipping the first 64 entries and starting on line 64,
//    it's easy to read the EBCDIC character as the line # and
//    the constant as its ASCII equivalent

// 🟦 Lookup tables

// 🔥 ChatGPT gave us this, for good or bad

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
	0x20,
	0xa0,
	0xe2,
	0xe4,
	0xe0,
	0xe1,
	0xe3,
	0xe5,
	0xe7,
	0xf1,
	0xa2,
	0x2e,
	0x3c,
	0x28,
	0x2b,
	0x7c,
	0x26,
	0xe9,
	0xeA,
	0xeB,
	0xe8,
	0xeD,
	0xeE,
	0xeF,
	0xeC,
	0xdF,
	0x21,
	0x24,
	0x2a,
	0x29,
	0x3b,
	0xaC,
	0x2d,
	0x2f,
	0xc2,
	0xc4,
	0xc0,
	0xc1,
	0xc3,
	0xc5,
	0xc7,
	0xd1,
	0xa6,
	0x2c,
	0x25,
	0x5f,
	0x3e,
	0x3f,
	0xf8,
	0xc9,
	0xcA,
	0xcB,
	0xc8,
	0xcD,
	0xcE,
	0xcF,
	0xcC,
	0x60,
	0x3a,
	0x23,
	0x40,
	0x27,
	0x3d,
	0x22,
	0xd8,
	0x61,
	0x62,
	0x63,
	0x64,
	0x65,
	0x66,
	0x67,
	0x68,
	0x69,
	0xaB,
	0xbB,
	0xf0,
	0xfD,
	0xfE,
	0xb1,
	0xb0,
	0x6a,
	0x6b,
	0x6c,
	0x6d,
	0x6e,
	0x6f,
	0x70,
	0x71,
	0x72,
	0xaA,
	0xbA,
	0xe6,
	0xb8,
	0xc6,
	0xa4,
	0xb5,
	0x7e,
	0x73,
	0x74,
	0x75,
	0x76,
	0x77,
	0x78,
	0x79,
	0x7a,
	0xa1,
	0xbF,
	0xd0,
	0xdD,
	0xdE,
	0xaE,
	0x5e,
	0xa3,
	0xa5,
	0xb7,
	0xa9,
	0xa7,
	0xb6,
	0xbC,
	0xbD,
	0xbE,
	0x5b,
	0x5d,
	0xaF,
	0xa8,
	0xb4,
	0xd7,
	0x7b,
	0x41,
	0x42,
	0x43,
	0x44,
	0x45,
	0x46,
	0x47,
	0x48,
	0x49,
	0xaD,
	0xf4,
	0xf6,
	0xf2,
	0xf3,
	0xf5,
	0x7d,
	0x4a,
	0x4b,
	0x4c,
	0x4d,
	0x4e,
	0x4f,
	0x50,
	0x51,
	0x52,
	0xb9,
	0xfB,
	0xfC,
	0xf9,
	0xfA,
	0xfF,
	0x5c,
	0xf7,
	0x53,
	0x54,
	0x55,
	0x56,
	0x57,
	0x58,
	0x59,
	0x5a,
	0xb2,
	0xd4,
	0xd6,
	0xd2,
	0xd3,
	0xd5,
	0x30,
	0x31,
	0x32,
	0x33,
	0x34,
	0x35,
	0x36,
	0x37,
	0x38,
	0x39,
	0xb3,
	0xdB,
	0xdC,
	0xd9,
	0xdA,
	0x9f,
}

// 🟦 Public functions

func E2A(e byte) byte {
	if e >= 64 {
		return EBCDIC[e-64]
	} else {
		return ' '
	}
}

func E2As(str string) string {
	ebcdic := []byte(str)
	ascii := make([]byte, len(ebcdic))
	for ix, char := range ebcdic {
		ascii[ix] = E2A(char)
	}
	return string(ascii)
}

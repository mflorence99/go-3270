package conv

import (
	"go3270/emulator/consts"
	"go3270/emulator/conv/cps"
)

// 🟧 EBCDIC -> Rune conversion

// 🟦 Lookup tables

var CPs = map[consts.LCID][]rune{
	0x00: cps.CP037,
	0xF1: cps.CP310,
}

// 🟦 Public functions

func E2Rune(lcid consts.LCID, e byte) rune {
	if e >= 64 {
		return CPs[lcid][e-64]
	} else {
		return '\u0020'
	}
}

func E2Runes(lcid consts.LCID, str string) string {
	ebcdic := []byte(str)
	runes := make([]rune, len(ebcdic))
	for ix, char := range ebcdic {
		runes[ix] = E2Rune(lcid, char)
	}
	return string(runes)
}

package conv

import "go3270/emulator/conv/cps"

var CPs = map[byte][]rune{
	0x00: cps.CP037,
	0xF1: cps.CP310,
}

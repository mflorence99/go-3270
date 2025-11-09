package screenshots

import (
	"go3270/emulator/conv"
	"go3270/emulator/utils"
)

// 🟧 Prefabricated outbound data streams for testing and performance measurement, captured via WireShark

var Index map[string][]byte

func init() {
	Index = make(map[string][]byte)
	Index["ge"] = utils.Flatten2Bytes(GE, conv.A2Es)
	Index["outline"] = utils.Flatten2Bytes(OUTLINE, conv.A2Es)
	Index["rpf-menu"] = RPF_MENU
	Index["symset0"] = SYMSET0
	Index["symset1"] = SYMSET1
	Index["termtest"] = TERMTEST
	Index["termtest-help"] = TERMTEST_HELP
}

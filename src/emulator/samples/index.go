package samples

import (
	"emulator/conv"
	"emulator/utils"
)

// 🟧 Prefabricated outbound data streams for testing
//    and performance measurement, captured via WireShark

var Index map[string][]byte

func init() {
	Index = make(map[string][]byte)
	Index["ge"] = utils.Flatten2Bytes(GE, conv.A2Es)
	Index["outline"] = utils.Flatten2Bytes(OUTLINE, conv.A2Es)
	Index["rm"] = utils.Flatten2Bytes(RM, conv.A2Es)
	Index["rpf-menu"] = RPF_MENU
	Index["short"] = utils.Flatten2Bytes(SHORT, conv.A2Es)
	Index["symset0"] = SYMSET0
	Index["symset1"] = SYMSET1
	Index["termtest"] = TERMTEST
	Index["termtest-help"] = TERMTEST_HELP
}

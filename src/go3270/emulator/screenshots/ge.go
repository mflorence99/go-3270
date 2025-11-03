package screenshots

import (
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/utils"
	"go3270/emulator/wcc"
)

// 🟧 Test GE order

var (
	raw = []any{
		consts.EW,
		wcc.WCC{}.Byte(),
		consts.SBA,
		conv.Addr2Bytes(0),
		"123-->GE",
		consts.GE,
		"GE<--456",
	}

	GE []byte
)

func init() {
	GE = utils.Flatten2Bytes(raw, conv.A2Es)
}

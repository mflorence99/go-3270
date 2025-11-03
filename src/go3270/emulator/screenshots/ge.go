package screenshots

import (
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/wcc"
)

// 🟧 Test GE order

var GE = []any{
	consts.EW,
	wcc.WCC{Alarm: true}.Byte(),
	consts.SBA,
	conv.Addr2Bytes(0),
	"123-->GE",
	consts.GE,
	"GE<--456",
}

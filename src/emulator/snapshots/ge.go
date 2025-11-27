//go:build dev

package snapshots

import (
	"emulator/conv"
	"emulator/types"
)

// 🟧 Test GE order

var GE = []any{
	types.EW,
	types.WCC{Alarm: true}.Bits(),
	types.SBA,
	conv.Addr2Bytes(0),
	"123-->GE",
	types.GE,
	"GE<--456",
}

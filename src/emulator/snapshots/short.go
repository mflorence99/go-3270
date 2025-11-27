//go:build dev

package snapshots

import (
	"emulator/conv"
	"emulator/types"
)

// 🟧 Test minimal page

var SHORT = []any{types.EW,
	types.WCC{Alarm: true}.Bits(),
	types.SBA,
	conv.Addr2Bytes(0),
	"ABC",
}

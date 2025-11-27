//go:build dev

package samples

import (
	"emulator/conv"
	"emulator/types"
)

// 🟧 Test RM with character attributes attribute

var RM = []any{
	types.EW,
	types.WCC{Reset: true}.Bits(),

	types.SBA,
	conv.Addr2Bytes(0),
	types.SFE,
	0x02,
	types.BASIC,
	(&types.Attrs{Protected: true}).Bits(),
	types.HIGHLIGHT,
	types.INTENSIFY,
	"Please enter something:",

	types.SBA,
	conv.Addr2Bytes(1),
	types.SA,
	types.COLOR,
	types.RED,

	types.SBA,
	conv.Addr2Bytes(2),
	types.SA,
	types.COLOR,
	types.PINK,

	types.SBA,
	conv.Addr2Bytes(3),
	types.SA,
	types.COLOR,
	types.BLUE,

	types.SBA,
	conv.Addr2Bytes(24),
	types.SFE,
	0x02,
	types.BASIC,
	(&types.Attrs{Protected: false}).Bits(),
	types.HIGHLIGHT,
	types.UNDERSCORE,
	"                  ",

	types.SBA,
	conv.Addr2Bytes(25),
	types.SA,
	types.COLOR,
	types.RED,

	types.SBA,
	conv.Addr2Bytes(26),
	types.SA,
	types.COLOR,
	types.PINK,

	types.SBA,
	conv.Addr2Bytes(27),
	types.SA,
	types.COLOR,
	types.BLUE,

	types.SBA,
	conv.Addr2Bytes(40),
	types.SF,
	(&types.Attrs{Autoskip: true, Numeric: true, Protected: true}).Bits(),
	"",

	types.SBA,
	conv.Addr2Bytes(25),
	types.IC,
}

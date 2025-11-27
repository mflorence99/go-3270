//go:build dev

package samples

import (
	"emulator/conv"
	"emulator/types"
)

// 🟧 Test outline attribute

var OUTLINE = []any{
	types.EW,
	types.WCC{Reset: true}.Bits(),

	types.SBA,
	conv.Addr2Bytes(0),
	types.SFE,
	0x03,
	types.BASIC,
	(&types.Attrs{Protected: true}).Bits(),
	types.HIGHLIGHT,
	types.INTENSIFY,
	types.OUTLINE,
	types.OUTLINE_BOTTOM | types.OUTLINE_RIGHT,
	"Header 1",

	types.SBA,
	conv.Addr2Bytes(10),
	types.SFE,
	0x03,
	types.BASIC,
	(&types.Attrs{Protected: true}).Bits(),
	types.HIGHLIGHT,
	types.INTENSIFY,
	types.OUTLINE,
	types.OUTLINE_BOTTOM | types.OUTLINE_LEFT,
	"Header 2",
	types.SF,
	(&types.Attrs{Protected: true}).Bits(),

	types.SBA,
	conv.Addr2Bytes(80),
	types.SFE,
	0x01,
	types.OUTLINE,
	types.OUTLINE_TOP | types.OUTLINE_RIGHT,
	"Cell 1/1",

	types.SBA,
	conv.Addr2Bytes(90),
	types.SFE,
	0x01,
	types.OUTLINE,
	types.OUTLINE_TOP | types.OUTLINE_LEFT,
	"Cell 1/2",
	types.SF,
	(&types.Attrs{Protected: true}).Bits(),

	types.SBA,
	conv.Addr2Bytes(160),
	types.SFE,
	0x01,
	types.OUTLINE,
	types.OUTLINE_TOP | types.OUTLINE_RIGHT,
	"Cell 2/1",

	types.SBA,
	conv.Addr2Bytes(170),
	types.SFE,
	0x01,
	types.OUTLINE,
	types.OUTLINE_TOP | types.OUTLINE_LEFT,
	"Cell 2/2",
	types.SF,
	(&types.Attrs{Protected: true}).Bits(),

	types.SBA,
	conv.Addr2Bytes(240),
	types.SFE,
	0x01,
	types.OUTLINE,
	types.OUTLINE_TOP | types.OUTLINE_RIGHT,
	"Cell 3/1",

	types.SBA,
	conv.Addr2Bytes(250),
	types.SFE,
	0x01,
	types.OUTLINE,
	types.OUTLINE_TOP | types.OUTLINE_LEFT,
	"Cell 3/2",
	types.SF,
	(&types.Attrs{Protected: true}).Bits(),

	types.SBA,
	conv.Addr2Bytes(240),
	types.SF,
	(&types.Attrs{Protected: true}).Bits(),
	"Ooops, I wanted data here",

	types.SBA,
	conv.Addr2Bytes(1),
	types.IC,
}

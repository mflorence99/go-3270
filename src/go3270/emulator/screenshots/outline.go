package screenshots

import (
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/wcc"
)

// 🟧 Test outline attribute

var OUTLINE = []any{
	consts.EW,
	wcc.WCC{Reset: true}.Byte(),

	consts.SBA,
	conv.Addr2Bytes(0),
	consts.SFE,
	0x03,
	consts.BASIC,
	(&attrs.Attrs{Protected: true}).Byte(),
	consts.HIGHLIGHT,
	consts.INTENSIFY,
	consts.OUTLINE,
	consts.OUTLINE_BOTTOM | consts.OUTLINE_RIGHT,
	"Header 1",

	consts.SBA,
	conv.Addr2Bytes(10),
	consts.SFE,
	0x03,
	consts.BASIC,
	(&attrs.Attrs{Protected: true}).Byte(),
	consts.HIGHLIGHT,
	consts.INTENSIFY,
	consts.OUTLINE,
	consts.OUTLINE_BOTTOM | consts.OUTLINE_LEFT,
	"Header 2",
	consts.SF,
	(&attrs.Attrs{Protected: true}).Byte(),

	consts.SBA,
	conv.Addr2Bytes(80),
	consts.SFE,
	0x01,
	consts.OUTLINE,
	consts.OUTLINE_TOP | consts.OUTLINE_RIGHT,
	"Cell 1/1",

	consts.SBA,
	conv.Addr2Bytes(90),
	consts.SFE,
	0x01,
	consts.OUTLINE,
	consts.OUTLINE_TOP | consts.OUTLINE_LEFT,
	"Cell 1/2",
	consts.SF,
	(&attrs.Attrs{Protected: true}).Byte(),

	consts.SBA,
	conv.Addr2Bytes(160),
	consts.SFE,
	0x01,
	consts.OUTLINE,
	consts.OUTLINE_TOP | consts.OUTLINE_RIGHT,
	"Cell 2/1",

	consts.SBA,
	conv.Addr2Bytes(170),
	consts.SFE,
	0x01,
	consts.OUTLINE,
	consts.OUTLINE_TOP | consts.OUTLINE_LEFT,
	"Cell 2/2",
	consts.SF,
	(&attrs.Attrs{Protected: true}).Byte(),

	consts.SBA,
	conv.Addr2Bytes(240),
	consts.SFE,
	0x01,
	consts.OUTLINE,
	consts.OUTLINE_TOP | consts.OUTLINE_RIGHT,
	"Cell 3/1",

	consts.SBA,
	conv.Addr2Bytes(250),
	consts.SFE,
	0x01,
	consts.OUTLINE,
	consts.OUTLINE_TOP | consts.OUTLINE_LEFT,
	"Cell 3/2",
	consts.SF,
	(&attrs.Attrs{Protected: true}).Byte(),

	consts.SBA,
	conv.Addr2Bytes(240),
	consts.SF,
	(&attrs.Attrs{Protected: true}).Byte(),
	"Ooops, I wanted data here",

	consts.SBA,
	conv.Addr2Bytes(1),
	consts.IC,
}

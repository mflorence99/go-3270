package qr

import (
	"go3270/emulator/stream"
	"go3270/emulator/types"
	"go3270/emulator/utils"
)

// 🟧 Query Reply structured field

type ColorSupport struct {
	SFID  types.SFID
	QCode types.QCode
	Flags byte
	NP    byte
	CAVs  []byte
}

// 🟦 Constructor

func NewColorSupport(monochrome bool) ColorSupport {
	cavs := make([]byte, 0)
	cavs = append(cavs, []byte{0x00, 0xF4}...)
	// TODO 🔥 somehow TSO gets confused when we add white 0xFF!
	// maybe the [0xFE, 0xFF] sequence is confused with the frame LT?
	// so I'm just pretending we only support 15 colors not 16
	// it doesn't seem to be a factor
	for ix := 1; ix < 15; ix++ {
		cavs = append(cavs, []byte{byte(ix + 240), utils.Ternary(monochrome, 0x00, byte(ix+240))}...)
	}
	return ColorSupport{
		SFID:  types.QUERY_REPLY,
		QCode: types.COLOR_SUPPORT,
		// 👇 flags appropriate for "not a printer"
		Flags: 0x00,
		NP:    byte(len(cavs) / 2),
		CAVs:  cavs,
	}
}

// 🟦 Public emitter function

func (s ColorSupport) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.Flags)
	chars = append(chars, s.NP)
	chars = append(chars, s.CAVs...)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
	"go3270/emulator/utils"
)

type ColorSupport struct {
	SFID  consts.SFID
	QCode consts.QCode
	Flags byte
	NP    byte
	CAVs  []byte
}

func NewColorSupport(monochrome bool) ColorSupport {
	cavs := make([]byte, 0)
	cavs = append(cavs, []byte{0x00, 0xF4}...)
	// TODO 🔥 somehow TSO gets confused when we add white 0xFF, maybe the 0xFe, 0xFF sequence is confused with the frame LT?
	for ix := 1; ix < 15; ix++ {
		cavs = append(cavs, []byte{byte(ix + 240), utils.Ternary(monochrome, 0x00, byte(ix+240))}...)
	}
	return ColorSupport{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.COLOR_SUPPORT,
		// 👇 flags appropriate for "not a printer"
		Flags: 0x00,
		NP:    byte(len(cavs) / 2),
		CAVs:  cavs,
	}
}

func (s ColorSupport) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	bytes = append(bytes, s.Flags)
	bytes = append(bytes, s.NP)
	bytes = append(bytes, s.CAVs...)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

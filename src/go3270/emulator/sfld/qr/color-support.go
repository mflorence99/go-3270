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
	CAVs  [][]byte
}

// 🔥 we just support the basic 7 colors, aliassing "black" to the default green, and we left the CLUT sort out what color is actually displayed

func NewColorSupport(monochrome bool) ColorSupport {
	cavs := make([][]byte, 16)
	cavs[0] = []byte{0x00, 0xF4}
	for ix := 1; ix < 16; ix++ {
		cavs[ix] = []byte{byte(ix + 240), utils.Ternary(monochrome, 0x00, byte(ix+240))}
	}
	return ColorSupport{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.COLOR_SUPPORT,
		// 👇 flags appropriate for "not a printer"
		Flags: 0x00,
		NP:    byte(len(cavs)),
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
	for _, cav := range s.CAVs {
		bytes = append(bytes, cav...)
	}

	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

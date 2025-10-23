package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"

	"golang.org/x/exp/maps"
)

type ColorSupport struct {
	SFID  consts.SFID
	QCode consts.QCode
	Flags byte
	NP    byte
	CAVs  [][]byte
}

func NewColorSupport(clut map[consts.Color][2]string) ColorSupport {
	colors := maps.Keys(clut)
	cavs := make([][]byte, len(colors))
	for ix, color := range colors {
		if color == consts.BLACK {
			// TODO 🔥 approximation to spec: 0x00 color is displayed as "green". In reality, we never use 0x00, but instead supplied defaults
			cavs[ix] = []byte{0x00, 0xF4}
		} else {
			cavs[ix] = []byte{byte(color), byte(color)}
		}
	}
	return ColorSupport{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.COLOR_SUPPORT,
		// 👇 flags appropriate for "not a printer"
		Flags: 0x00,
		NP:    byte(len(clut)),
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

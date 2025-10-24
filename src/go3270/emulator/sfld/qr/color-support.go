package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type ColorSupport struct {
	SFID  consts.SFID
	QCode consts.QCode
	Flags byte
	NP    byte
	CAVs  [][]byte
}

// 🔥 we just support every color, aliassing "black" to the default green, and we left the CLUT sort out what color is actually displayed

func NewColorSupport() ColorSupport {
	cavs := make([][]byte, 16)
	cavs[0] = []byte{0x00, 0xF4}
	cavs[1] = []byte{0xF1, 0xF1}
	cavs[2] = []byte{0xF2, 0xF2}
	cavs[3] = []byte{0xF3, 0xF3}
	cavs[4] = []byte{0xF4, 0xF4}
	cavs[5] = []byte{0xF5, 0xF5}
	cavs[6] = []byte{0xF6, 0xF6}
	cavs[7] = []byte{0xF7, 0xF7}
	cavs[8] = []byte{0xF8, 0xF8}
	cavs[9] = []byte{0xF9, 0xF9}
	cavs[10] = []byte{0xFA, 0xFA}
	cavs[11] = []byte{0xFB, 0xFB}
	cavs[12] = []byte{0xFC, 0xFC}
	cavs[13] = []byte{0xFD, 0xFD}
	cavs[14] = []byte{0xFE, 0xFE}
	cavs[15] = []byte{0xFF, 0xFF}
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

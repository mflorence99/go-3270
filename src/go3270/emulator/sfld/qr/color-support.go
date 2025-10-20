package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type ColorSupport struct {
	SFID  consts.SFID
	QCode consts.QCode
	CLUT  map[consts.Color][2]string
}

func NewColorSupport(clut map[consts.Color][2]string) ColorSupport {
	return ColorSupport{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.COLOR_SUPPORT,
		CLUT:  clut,
	}
}

func (s ColorSupport) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags appropriate for "not a printer"
	bytes = append(bytes, 0b00000000)
	// 👇 extract color support from the CLUT
	bytes = append(bytes, byte(len(s.CLUT)))
	for k := range s.CLUT {
		if k == 0xF0 {
			// 🔥 approximation to spec: 0x00 color is displayed as "green". In reality, we never use 0x00, but instead supplied defaults
			bytes = append(bytes, 0x00)
			bytes = append(bytes, byte(consts.GREEN))
		} else {
			// 👇 normal case: requested to actual color mapping
			bytes = append(bytes, byte(k))
			bytes = append(bytes, byte(k))
		}
	}
	return bytes, uint16(len(bytes) + 2)
}

func (s ColorSupport) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

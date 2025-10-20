package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type Highlighting struct {
	SFID  consts.SFID
	QCode consts.QCode
}

func NewHighlighting() Highlighting {
	return Highlighting{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.HIGHLIGHTING,
	}
}

func (s Highlighting) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 indicate full support for all highlighting
	bytes = append(bytes, 0x05)
	bytes = append(bytes, []byte{0x00, byte(consts.NOHILITE)}...)
	bytes = append(bytes, []byte{0xF1, byte(consts.BLINK)}...)
	bytes = append(bytes, []byte{0xF2, byte(consts.REVERSE)}...)
	bytes = append(bytes, []byte{0xF4, byte(consts.UNDERSCORE)}...)
	bytes = append(bytes, []byte{0xF5, byte(consts.INTENSITY)}...)
	return bytes, uint16(len(bytes) + 2)
}

func (s Highlighting) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

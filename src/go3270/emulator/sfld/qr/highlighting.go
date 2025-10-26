package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type Highlighting struct {
	SFID  consts.SFID
	QCode consts.QCode
	NP    byte
	HAVs  [][]byte
}

func NewHighlighting() Highlighting {
	havs := make([][]byte, 5)
	havs[0] = []byte{0x00, byte(consts.NOHILITE)}
	havs[1] = []byte{0xF1, byte(consts.BLINK)}
	havs[2] = []byte{0xF2, byte(consts.REVERSE)}
	havs[3] = []byte{0xF3, byte(consts.UNDERSCORE)}
	havs[4] = []byte{0xF4, byte(consts.HIGHLIGHT)}
	return Highlighting{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.HIGHLIGHTING,
		NP:    0x05,
		HAVs:  havs,
	}
}

func (s Highlighting) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	bytes = append(bytes, s.NP)
	for _, hav := range s.HAVs {
		bytes = append(bytes, hav...)
	}
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

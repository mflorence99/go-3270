package qr

import (
	"go3270/emulator/stream"
	"go3270/emulator/types"
)

// 🟧 Query Reply structured field

type Highlighting struct {
	SFID  types.SFID
	QCode types.QCode
	NP    byte
	HAVs  [][]byte
}

// 🟦 Constructor

func NewHighlighting() Highlighting {
	havs := make([][]byte, 5)
	havs[0] = []byte{0x00, byte(types.NO_HILITE)}
	havs[1] = []byte{0xF1, byte(types.BLINK)}
	havs[2] = []byte{0xF2, byte(types.REVERSE)}
	havs[3] = []byte{0xF4, byte(types.UNDERSCORE)}
	havs[4] = []byte{0xF8, byte(types.INTENSIFY)}
	return Highlighting{
		SFID:  types.QUERY_REPLY,
		QCode: types.HIGHLIGHTING,
		NP:    0x05,
		HAVs:  havs,
	}
}

// 🟦 Public emitter function

func (s Highlighting) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.NP)
	for _, hav := range s.HAVs {
		chars = append(chars, hav...)
	}
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

package qr

import (
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Highlighting) pp 6-65 to 6-66

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
	havs[1] = []byte{0xf1, byte(types.BLINK)}
	havs[2] = []byte{0xf2, byte(types.REVERSE)}
	havs[3] = []byte{0xf4, byte(types.UNDERSCORE)}
	havs[4] = []byte{0xf8, byte(types.INTENSIFY)}
	return Highlighting{
		SFID:  types.QUERY_REPLY,
		QCode: types.HIGHLIGHTING,
		NP:    0x05,
		HAVs:  havs,
	}
}

// 🟦 Public emitter function

func (s Highlighting) Put(in iface.Inbound) {
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

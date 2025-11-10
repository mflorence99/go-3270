package qr

import (
	"go3270/emulator/stream"
	"go3270/emulator/types"
)

// 🟧 Query Reply structured field

type FieldValidation struct {
	SFID  types.SFID
	QCode types.QCode
	Types byte
}

// 🟦 Constructor

func NewFieldValidation() FieldValidation {
	return FieldValidation{
		SFID:  types.QUERY_REPLY,
		QCode: types.FIELD_VALIDATION,
		// 👇 we support mandatory fill and entry, plus trigger
		Types: 0b00000111,
	}
}

// 🟦 Public emitter function

func (s FieldValidation) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	chars = append(chars, s.Types)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

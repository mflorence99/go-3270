package qr

import (
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Field Validation) p 6-59

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

func (s FieldValidation) Put(in iface.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	chars = append(chars, s.Types)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type FieldValidation struct {
	SFID  consts.SFID
	QCode consts.QCode
	Types byte
}

func NewFieldValidation() FieldValidation {
	return FieldValidation{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.FIELD_VALIDATION,
		// 👇 we support mandatory fill and entry, plus trigger
		Types: 0b00000111,
	}
}

func (s FieldValidation) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	bytes = append(bytes, s.Types)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

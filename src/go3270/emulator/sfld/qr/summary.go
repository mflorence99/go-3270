package qr

import (
	"go3270/emulator/stream"
	"go3270/emulator/types"
)

// 🟧 Query Reply structured field

type Summary struct {
	SFID  types.SFID
	QCode types.QCode
	List  []types.QCode
}

// 🟦 Constructor

func NewSummary(list []types.QCode) Summary {
	return Summary{
		SFID:  types.QUERY_REPLY,
		QCode: types.SUMMARY,
		List:  list,
	}
}

// 🟦 Public emitter function

func (s Summary) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	for _, qcode := range s.List {
		chars = append(chars, byte(qcode))
	}
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

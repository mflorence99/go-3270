package qr

import (
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Summary) p 6-96

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

func (s Summary) Put(in iface.Inbound) {
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

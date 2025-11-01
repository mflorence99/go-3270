package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type Summary struct {
	SFID  consts.SFID
	QCode consts.QCode
	List  []consts.QCode
}

func NewSummary(list []consts.QCode) Summary {
	return Summary{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.SUMMARY,
		List:  list,
	}
}

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

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
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	for qcode := range s.List {
		bytes = append(bytes, byte(qcode))
	}
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

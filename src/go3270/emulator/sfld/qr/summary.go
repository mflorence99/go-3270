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

func (s Summary) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	for qcode := range s.List {
		bytes = append(bytes, byte(qcode))
	}
	return bytes, uint16(len(bytes) + 2)
}

func (s Summary) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

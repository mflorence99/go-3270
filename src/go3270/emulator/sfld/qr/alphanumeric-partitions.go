package qr

import (
	"encoding/binary"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type AlphanumericPartitions struct {
	SFID  consts.SFID
	QCode consts.QCode
	Cols  int
	Rows  int
}

func NewAlphanumericPartitions(cols, rows int) AlphanumericPartitions {
	return AlphanumericPartitions{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.SUMMARY,
		Cols:  cols,
		Rows:  rows,
	}
}

func (s AlphanumericPartitions) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 max number of partitions
	bytes = append(bytes, 0x00)
	// 👇 total storage
	total := uint16(s.Cols * s.Rows)
	binary.LittleEndian.AppendUint16(bytes, uint16(total))
	return bytes, uint16(len(bytes) + 2)
}

func (s AlphanumericPartitions) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

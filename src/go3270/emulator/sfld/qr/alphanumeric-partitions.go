package qr

import (
	"encoding/binary"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type AlphanumericPartitions struct {
	SFID  consts.SFID
	QCode consts.QCode
	NA    byte
	M     uint16
	Flags byte
}

func NewAlphanumericPartitions(cols, rows int) AlphanumericPartitions {
	return AlphanumericPartitions{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.ALPHANUMERIC_PARTITIONS,
		NA:    0x00,
		M:     uint16(cols * rows),
		Flags: 0x00,
	}
}

func (s AlphanumericPartitions) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	bytes = append(bytes, s.NA)
	bytes = binary.BigEndian.AppendUint16(bytes, s.M)
	bytes = append(bytes, s.Flags)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

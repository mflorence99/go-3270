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
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.NA)
	chars = binary.BigEndian.AppendUint16(chars, s.M)
	chars = append(chars, s.Flags)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

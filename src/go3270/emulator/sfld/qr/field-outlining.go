package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type FieldOutlining struct {
	SFID  consts.SFID
	QCode consts.QCode
	Flag  byte
	Sep   byte
	VPOS  byte
	HPOS  byte
	HPOS0 byte
	HPOS1 byte
}

func NewFieldOutlining() FieldOutlining {
	return FieldOutlining{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.FIELD_OUTLINING,
		// 👇 fill as best we can for a non-printer
		Flag:  0x00,
		Sep:   0b10000000,
		VPOS:  0x00,
		HPOS:  0x00,
		HPOS0: 0x00,
		HPOS1: 0x00,
	}
}

func (s FieldOutlining) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	bytes = append(bytes, s.Flag)
	bytes = append(bytes, s.Sep)
	bytes = append(bytes, s.VPOS)
	bytes = append(bytes, s.HPOS)
	bytes = append(bytes, s.HPOS0)
	bytes = append(bytes, s.HPOS1)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

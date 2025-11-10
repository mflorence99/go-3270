package qr

import (
	"go3270/emulator/stream"
	"go3270/emulator/types"
)

// 🟧 Query Reply structured field

type FieldOutlining struct {
	SFID  types.SFID
	QCode types.QCode
	Flag  byte
	Sep   byte
	VPOS  byte
	HPOS  byte
	HPOS0 byte
	HPOS1 byte
}

// 🟦 Constructor

func NewFieldOutlining() FieldOutlining {
	return FieldOutlining{
		SFID:  types.QUERY_REPLY,
		QCode: types.FIELD_OUTLINING,
		// 👇 fill as best we can for a non-printer
		Flag:  0x00,
		Sep:   0b10000000,
		VPOS:  0x00,
		HPOS:  0x00,
		HPOS0: 0x00,
		HPOS1: 0x00,
	}
}

// 🟦 Public emitter function

func (s FieldOutlining) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	chars = append(chars, s.Flag)
	chars = append(chars, s.Sep)
	chars = append(chars, s.VPOS)
	chars = append(chars, s.HPOS)
	chars = append(chars, s.HPOS0)
	chars = append(chars, s.HPOS1)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

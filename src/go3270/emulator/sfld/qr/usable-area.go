package qr

import (
	"encoding/binary"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type UsableArea struct {
	SFID   consts.SFID
	QCode  consts.QCode
	Flags1 byte
	Flags2 byte
	W      uint16
	H      uint16
}

func NewUsableArea(cols, rows int) UsableArea {
	return UsableArea{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.USABLE_AREA,
		// 👇 12/14 bit addressing
		Flags1: 0b00000001,
		// 👇 dimensions in cells (not pells)
		Flags2: 0b00000000,
		W:      uint16(cols),
		H:      uint16(rows),
	}
}

func (s UsableArea) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.Flags1)
	chars = append(chars, s.Flags2)
	chars = binary.BigEndian.AppendUint16(chars, s.W)
	chars = binary.BigEndian.AppendUint16(chars, s.H)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

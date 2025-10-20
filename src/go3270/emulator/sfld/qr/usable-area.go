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

func (s UsableArea) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	bytes = append(bytes, s.Flags1)
	bytes = append(bytes, s.Flags2)
	bytes = binary.BigEndian.AppendUint16(bytes, s.W)
	bytes = binary.BigEndian.AppendUint16(bytes, s.H)
	return bytes, uint16(len(bytes) + 2)
}

func (s UsableArea) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

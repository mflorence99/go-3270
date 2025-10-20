package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type UsableArea struct {
	SFID  consts.SFID
	QCode consts.QCode
	Cols  int
	Rows  int
}

func NewUsableArea(cols, rows int) UsableArea {
	return UsableArea{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.USABLE_AREA,
		Cols:  cols,
		Rows:  rows,
	}
}

func (s UsableArea) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 12/14 bit addressing
	bytes = append(bytes, 0b00000001)
	// 👇 dimensions in cells (not pells)
	bytes = append(bytes, 0b00000000)
	// 👇 width
	bytes = append(bytes, 0b00000000)
	bytes = append(bytes, byte(s.Cols))
	// 👇 height
	bytes = append(bytes, 0b00000000)
	bytes = append(bytes, byte(s.Rows))
	// 🔥 TBD
	return bytes, uint16(len(bytes) + 2)
}

func (s UsableArea) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

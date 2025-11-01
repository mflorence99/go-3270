package qr

import (
	"encoding/binary"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type ImplicitPartition struct {
	SFID   consts.SFID
	QCode  consts.QCode
	Flags1 []byte
	L      byte
	SDPID  byte
	Flags2 byte
	WD     uint16
	HD     uint16
	WA     uint16
	HA     uint16
}

func NewImplicitPartition(cols, rows int) ImplicitPartition {
	return ImplicitPartition{
		SFID:   consts.QUERY_REPLY,
		QCode:  consts.IMPLICIT_PARTITION,
		Flags1: []byte{0x00, 0x00},
		L:      0x0B,
		SDPID:  0x01,
		Flags2: 0x00,
		WD:     uint16(cols),
		HD:     uint16(rows),
		WA:     uint16(cols),
		HA:     uint16(rows),
	}
}

func (s ImplicitPartition) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.Flags1...)
	chars = append(chars, s.L)
	chars = append(chars, s.SDPID)
	chars = append(chars, s.Flags2)
	chars = binary.BigEndian.AppendUint16(chars, s.WD)
	chars = binary.BigEndian.AppendUint16(chars, s.HD)
	chars = binary.BigEndian.AppendUint16(chars, s.WA)
	chars = binary.BigEndian.AppendUint16(chars, s.HA)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

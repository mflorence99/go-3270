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

func (s ImplicitPartition) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	bytes = append(bytes, s.Flags1...)
	bytes = append(bytes, s.L)
	bytes = append(bytes, s.SDPID)
	bytes = append(bytes, s.Flags2)
	bytes = binary.BigEndian.AppendUint16(bytes, s.WD)
	bytes = binary.BigEndian.AppendUint16(bytes, s.HD)
	bytes = binary.BigEndian.AppendUint16(bytes, s.WA)
	bytes = binary.BigEndian.AppendUint16(bytes, s.HA)
	return bytes, uint16(len(bytes) + 2)
}

func (s ImplicitPartition) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

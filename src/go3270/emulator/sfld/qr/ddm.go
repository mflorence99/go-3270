package qr

import (
	"encoding/binary"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type DDM struct {
	SFID   consts.SFID
	QCode  consts.QCode
	Flags  []byte
	LimIn  uint16
	LimOut uint16
	NSS    byte
	DDMSS  byte
}

func NewDDM() DDM {
	return DDM{
		SFID:   consts.QUERY_REPLY,
		QCode:  consts.DDM,
		Flags:  []byte{0x00, 0x00},
		LimIn:  uint16(4096),
		LimOut: uint16(4096),
		NSS:    1,
		DDMSS:  1,
	}
}

func (s DDM) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	bytes = append(bytes, s.Flags...)
	bytes = binary.BigEndian.AppendUint16(bytes, s.LimIn)
	bytes = binary.BigEndian.AppendUint16(bytes, s.LimOut)
	bytes = append(bytes, s.NSS)
	bytes = append(bytes, s.DDMSS)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

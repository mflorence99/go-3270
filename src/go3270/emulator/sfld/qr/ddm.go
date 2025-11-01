package qr

import (
	"encoding/binary"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

// 🟧 Query Reply structured field

type DDM struct {
	SFID   consts.SFID
	QCode  consts.QCode
	Flags  []byte
	LimIn  uint16
	LimOut uint16
	NSS    byte
	DDMSS  byte
}

// 🟦 Constructor

func NewDDM() DDM {
	return DDM{
		SFID:   consts.QUERY_REPLY,
		QCode:  consts.DDM,
		Flags:  []byte{0x00, 0x00},
		LimIn:  uint16(4096 * 4),
		LimOut: uint16(4096 * 4),
		NSS:    1,
		DDMSS:  1,
	}
}

// 🟦 Public emitter function

func (s DDM) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.Flags...)
	chars = binary.BigEndian.AppendUint16(chars, s.LimIn)
	chars = binary.BigEndian.AppendUint16(chars, s.LimOut)
	chars = append(chars, s.NSS)
	chars = append(chars, s.DDMSS)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

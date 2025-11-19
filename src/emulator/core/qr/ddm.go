package qr

import (
	"emulator/iface"
	"emulator/types"
	"encoding/binary"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Distributed Data Management) pp 6-52 to 6-54

type DDM struct {
	SFID   types.SFID
	QCode  types.QCode
	Flags  []byte
	LimIn  uint16
	LimOut uint16
	NSS    byte
	DDMSS  byte
}

// 🟦 Constructor

func NewDDM() DDM {
	return DDM{
		SFID:   types.QUERY_REPLY,
		QCode:  types.DDM,
		Flags:  []byte{0x00, 0x00},
		LimIn:  uint16(4096 * 4),
		LimOut: uint16(4096 * 4),
		NSS:    1,
		DDMSS:  1,
	}
}

// 🟦 Public emitter function

func (s DDM) Put(in iface.Inbound) {
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

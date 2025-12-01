package qr

import (
	"emulator/iface"
	"emulator/types"
	"encoding/binary"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Usable Area) pp 6-100 to 6-109

type UsableArea struct {
	SFID   types.SFID
	QCode  types.QCode
	Flags1 byte
	Flags2 byte
	W      uint16
	H      uint16
	Units  byte
	Xr     []byte
	Yr     []byte
	AW     byte
	AH     byte
	BuffSz uint16
}

// 🟦 Constructor

func NewUsableArea(cols, rows uint, fontWidth, fontHeight float64) UsableArea {
	return UsableArea{
		SFID:  types.QUERY_REPLY,
		QCode: types.USABLE_AREA,
		// 👇 12/14 bit addressing
		Flags1: 0b00000001,
		// 👇 dimensions in cells (not pells)
		Flags2: 0b00000000,
		W:      uint16(cols),
		H:      uint16(rows),
		// 👇 mm for some reason
		Units: 0x01,
		// 👇 magic numbers from x3270 implementation
		Xr: []byte{0x00, 0x00, 0x00, 0x00},
		Yr: []byte{0x00, 0x00, 0x00, 0x00},
		// 👇 like SDW/SDH in character-sets
		AW:     byte(fontWidth),
		AH:     byte(fontHeight),
		BuffSz: uint16(cols * rows),
	}
}

// 🟦 Public emitter function

func (s UsableArea) Put(in iface.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.Flags1)
	chars = append(chars, s.Flags2)
	chars = binary.BigEndian.AppendUint16(chars, s.W)
	chars = binary.BigEndian.AppendUint16(chars, s.H)
	chars = append(chars, s.Units)
	chars = append(chars, s.Xr...)
	chars = append(chars, s.Yr...)
	chars = append(chars, s.AW)
	chars = append(chars, s.AH)
	chars = binary.BigEndian.AppendUint16(chars, s.BuffSz)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

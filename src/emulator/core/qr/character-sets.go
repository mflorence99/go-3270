package qr

import (
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Character Sets) pp 6-28 to 6-35

type CharacterSets struct {
	SFID  types.SFID
	QCode types.QCode
	Flag1 byte
	Flag2 byte
	SDW   byte
	SDH   byte
	FORM  []byte
	DL    byte
	Descs []CharacterSetDesc
}

type CharacterSetDesc struct {
	SET  byte
	Flag byte
	LCID byte
}

// 🟦 Constructor

func NewCharacterSets(fontWidth, fontHeight float64) CharacterSets {
	return CharacterSets{
		SFID:  types.QUERY_REPLY,
		QCode: types.CHARACTER_SETS,
		Flag1: 0b10000010,
		Flag2: 0b00000000,
		SDW:   byte(fontWidth),
		SDH:   byte(fontHeight),
		FORM:  []byte{0x00, 0x00, 0x00, 0x00},
		// 🔥 we really want len(CharacterSetDesc{})
		// 👇 length of each char set, of which we support 2
		DL: 3,
		Descs: []CharacterSetDesc{
			{SET: 0x00, Flag: 0b00010000, LCID: 0x00},
			{SET: 0x01, Flag: 0b00000000, LCID: 0xf1},
		},
	}
}

// 🟦 Public emitter function

func (s CharacterSets) Put(in iface.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags indicating basic support
	chars = append(chars, s.Flag1)
	chars = append(chars, s.Flag2)
	chars = append(chars, s.SDW)
	chars = append(chars, s.SDH)
	chars = append(chars, s.FORM...)
	chars = append(chars, s.DL)
	for _, desc := range s.Descs {
		chars = append(chars, desc.SET)
		chars = append(chars, desc.Flag)
		chars = append(chars, desc.LCID)
	}
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

package qr

import (
	"emulator/iface"
	"emulator/types"
	"emulator/utils"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Color) pp 6-36 to 6-38

type ColorSupport struct {
	SFID  types.SFID
	QCode types.QCode
	Flags byte
	NP    byte
	CAVs  []byte
}

// 🟦 Constructor

func NewColorSupport(monochrome bool) ColorSupport {
	cavs := make([]byte, 0)
	cavs = append(cavs, []byte{0x00, 0xf4}...)
	for ix := 1; ix < 16; ix++ {
		cavs = append(cavs, []byte{byte(ix + 240), utils.Ternary(monochrome, 0x00, byte(ix+240))}...)
	}
	return ColorSupport{
		SFID:  types.QUERY_REPLY,
		QCode: types.COLOR_SUPPORT,
		// 👇 flags appropriate for "not a printer"
		Flags: 0x00,
		NP:    byte(len(cavs) / 2),
		CAVs:  cavs,
	}
}

// 🟦 Public emitter function

func (s ColorSupport) Put(in iface.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags and data
	chars = append(chars, s.Flags)
	chars = append(chars, s.NP)
	chars = append(chars, s.CAVs...)
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

package qr

import (
	"emulator/conv"
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (RPQ Names) pp 6-90 to 6-91

type RPQNames struct {
	SFID    types.SFID
	QCode   types.QCode
	Device  []byte
	Model   []byte
	RPQName string
}

// 🟦 Constructor

func NewRPQNames() RPQNames {
	return RPQNames{
		SFID:    types.QUERY_REPLY,
		QCode:   types.RPQ_NAMES,
		Device:  []byte{0x00, 0x00, 0x00, 0x00},
		Model:   []byte{0x00, 0x00, 0x00, 0x00},
		RPQName: "go3270",
	}
}

// 🟦 Public emitter function

func (s RPQNames) Put(in iface.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 RPQ data
	chars = append(chars, s.Device...)
	chars = append(chars, s.Model...)
	chars = append(chars, byte(len(s.RPQName)+1))
	for _, a := range s.RPQName {
		chars = append(chars, conv.A2E(byte(a)))
	}
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

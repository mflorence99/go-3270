package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/stream"
)

type RPQNames struct {
	SFID    consts.SFID
	QCode   consts.QCode
	Device  []byte
	Model   []byte
	RPQName string
}

func NewRPQNames() RPQNames {
	return RPQNames{
		SFID:    consts.QUERY_REPLY,
		QCode:   consts.RPQ_NAMES,
		Device:  []byte{0x00, 0x00, 0x00, 0x00},
		Model:   []byte{0x00, 0x00, 0x00, 0x00},
		RPQName: "go3270",
	}
}

func (s RPQNames) Put(in *stream.Inbound) {
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

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

func (s RPQNames) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 RPQ data
	bytes = append(bytes, s.Device...)
	bytes = append(bytes, s.Model...)
	bytes = append(bytes, byte(len(s.RPQName)+1))
	for _, a := range s.RPQName {
		bytes = append(bytes, conv.A2E(byte(a)))
	}
	return bytes, uint16(len(bytes) + 2)
}

func (s RPQNames) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

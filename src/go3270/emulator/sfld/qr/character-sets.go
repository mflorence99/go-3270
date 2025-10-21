package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type CharacterSets struct {
	SFID  consts.SFID
	QCode consts.QCode
	Flag1 byte
	Flag2 byte
}

func NewCharacterSets() CharacterSets {
	return CharacterSets{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.CHARACTER_SETS,
		Flag1: 0b00000000,
		Flag2: 0b00000000,
	}
}

func (s CharacterSets) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags indicating basic support
	bytes = append(bytes, s.Flag1)
	bytes = append(bytes, s.Flag2)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

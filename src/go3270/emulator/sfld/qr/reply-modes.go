package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

// 🟧 Query Reply structured field

type ReplyModes struct {
	SFID  consts.SFID
	QCode consts.QCode
	Modes []consts.Mode
}

// 🟦 Constructor

func NewReplyModes() ReplyModes {
	return ReplyModes{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.REPLY_MODES,
		Modes: []consts.Mode{
			consts.FIELD_MODE,
			consts.EXTENDED_FIELD_MODE,
			consts.CHARACTER_MODE,
		},
	}
}

// 🟦 Public emitter function

func (s ReplyModes) Put(in *stream.Inbound) {
	chars := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	for _, mode := range s.Modes {
		chars = append(chars, byte(mode))
	}
	in.Put16(uint16(len(chars) + 2))
	in.PutSlice(chars)
}

package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type ReplyModes struct {
	SFID  consts.SFID
	QCode consts.QCode
	Modes []consts.Mode
}

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

func (s ReplyModes) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	for _, mode := range s.Modes {
		bytes = append(bytes, byte(mode))
	}
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

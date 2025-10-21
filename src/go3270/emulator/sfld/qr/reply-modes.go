package qr

import (
	"go3270/emulator/consts"
	"go3270/emulator/stream"
)

type ReplyModes struct {
	SFID  consts.SFID
	QCode consts.QCode
	Modes []byte
}

func NewReplyModes() ReplyModes {
	return ReplyModes{
		SFID:  consts.QUERY_REPLY,
		QCode: consts.REPLY_MODES,
		// 👇 field, extended field and character (SF, SFE and SA)
		Modes: []byte{0x00, 0x01, 0x02},
	}
}

func (s ReplyModes) Put(in *stream.Inbound) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	bytes = append(bytes, s.Modes...)
	in.Put16(uint16(len(bytes) + 2))
	in.PutSlice(bytes)
}

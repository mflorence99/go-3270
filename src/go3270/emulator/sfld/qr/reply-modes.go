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
		Modes: []byte{0x00, 0x02, 0x02},
	}
}

func (s ReplyModes) Bytes() ([]byte, uint16) {
	bytes := []byte{
		byte(s.SFID),
		byte(s.QCode),
	}
	// 👇 flags
	bytes = append(bytes, s.Modes...)
	return bytes, uint16(len(bytes) + 2)
}

func (s ReplyModes) Put(in *stream.Inbound) {
	bytes, len := s.Bytes()
	in.Put16(len)
	in.PutSlice(bytes)
}

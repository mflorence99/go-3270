package qr

import (
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

type ReplyModes struct {
	SFID  types.SFID
	QCode types.QCode
	Modes []types.Mode
}

// 🟦 Constructor

func NewReplyModes() ReplyModes {
	return ReplyModes{
		SFID:  types.QUERY_REPLY,
		QCode: types.REPLY_MODES,
		Modes: []types.Mode{
			types.FIELD_MODE,
			types.EXTENDED_FIELD_MODE,
			types.CHARACTER_MODE,
		},
	}
}

// 🟦 Public emitter function

func (s ReplyModes) Put(in iface.Inbound) {
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

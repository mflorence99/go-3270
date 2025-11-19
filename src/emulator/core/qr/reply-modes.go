package qr

import (
	"emulator/iface"
	"emulator/types"
)

// 🟧 Query Reply structured field

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Query Reply (Reply Modes) pp 6-89 to 6-90

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

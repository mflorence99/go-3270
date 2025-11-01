package stream

import (
	"encoding/binary"
)

// 🟧 Inbound (3270 -> app) data stream

type Inbound struct {
	chars []byte
}

// 🟦 Constructor

func NewInbound() *Inbound {
	in := new(Inbound)
	in.chars = []byte{}
	return in
}

// 🟦 Public functions

func (in *Inbound) Bytes() []byte {
	return in.chars
}

func (in *Inbound) Put(char byte) []byte {
	in.chars = append(in.chars, char)
	return in.chars
}

func (in *Inbound) Put16(chars uint16) []byte {
	in.chars = binary.BigEndian.AppendUint16(in.chars, chars)
	return in.chars
}

func (in *Inbound) PutSlice(slice []byte) []byte {
	in.chars = append(in.chars, slice...)
	return in.chars
}

package stream

import (
	"encoding/binary"
)

// 🔥 "Inbound" data flows from the 3270 ie this code to the application

type Inbound struct {
	chars []byte
}

func NewInbound() *Inbound {
	in := new(Inbound)
	in.chars = []byte{}
	return in
}

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

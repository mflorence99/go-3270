package stream

import (
	"encoding/binary"
)

// 🔥 "Inbound" data flows from the 3270 ie this code to the application

type Inbound struct {
	bytes []byte
}

func NewInbound() *Inbound {
	in := new(Inbound)
	in.bytes = []byte{}
	return in
}

func (in *Inbound) Bytes() []byte {
	return in.bytes
}

func (in *Inbound) Put(char byte) []byte {
	in.bytes = append(in.bytes, char)
	return in.bytes
}

func (in *Inbound) Put16(chars uint16) []byte {
	in.bytes = binary.BigEndian.AppendUint16(in.bytes, chars)
	return in.bytes
}

func (in *Inbound) PutSlice(slice []byte) []byte {
	in.bytes = append(in.bytes, slice...)
	return in.bytes
}

package iface

type Inbound interface {
	Bytes() []byte
	Put(char byte) []byte
	Put16(chars uint16) []byte
	PutSlice(slice []byte) []byte
}

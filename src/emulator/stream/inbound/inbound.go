package inbound

// 🔥 "Inbound" data flows from the 3270 ie this code to the application

type Inbound struct {
	bytes []byte
}

func New() *Inbound {
	in := new(Inbound)
	in.bytes = []byte{}
	return in
}

func (in *Inbound) Put(u8 byte) []byte {
	in.bytes = append(in.bytes, u8)
	return in.bytes
}

func (in *Inbound) Put16(u16 uint16) []byte {
	in.bytes = append(in.bytes, byte(u16>>8))
	in.bytes = append(in.bytes, byte(u16&0x00ff))
	return in.bytes
}

func (in *Inbound) PutSlice(slice []byte) []byte {
	in.bytes = append(in.bytes, slice...)
	return in.bytes
}

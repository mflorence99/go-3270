package stream

// 🔥 "Inbound" data flows from the 3270 ie this code to the application

type Inbound struct {
	bytes []byte
}

func NewInbound() *Inbound {
	in := new(Inbound)
	in.bytes = []byte{}
	return in
}

func (in *Inbound) Put(char byte) []byte {
	in.bytes = append(in.bytes, char)
	return in.bytes
}

func (in *Inbound) Put16(chars uint16) []byte {
	in.bytes = append(in.bytes, byte(chars>>8))
	in.bytes = append(in.bytes, byte(chars&0x00ff))
	return in.bytes
}

func (in *Inbound) PutSlice(slice []byte) []byte {
	in.bytes = append(in.bytes, slice...)
	return in.bytes
}

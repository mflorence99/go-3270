package inbound

import (
	"go3270/emulator/pubsub"
)

type Producer struct {
	bus *pubsub.Bus
}

func NewProducer(bus *pubsub.Bus) *Producer {
	i := new(Producer)
	i.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	i.bus.SubClose(i.close)
	return i
}

func (i *Producer) close() {}

func (i *Producer) Produce() {
	bytes := make([]byte, 0)
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "palegreen",
		EBCDIC: true,
		Title:  "Inbound",
	}
	i.bus.PubDump(dmp)
	i.bus.PubInbound(bytes)
}

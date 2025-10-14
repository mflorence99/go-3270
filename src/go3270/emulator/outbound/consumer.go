package outbound

import (
	"go3270/emulator/pubsub"
)

type Consumer struct {
	bus *pubsub.Bus
}

func NewConsumer(bus *pubsub.Bus) *Consumer {
	o := new(Consumer)
	o.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	o.bus.SubClose(o.close)
	o.bus.SubOutbound(o.consume)
	return o
}

func (o *Consumer) close() {}

func (o *Consumer) consume(bytes []byte) {
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "yellow",
		EBCDIC: true,
		Title:  "Outbound",
	}
	o.bus.PubDump(dmp)
}

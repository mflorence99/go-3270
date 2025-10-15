package outbound

import (
	"go3270/emulator/pubsub"
)

type Consumer struct {
	bus *pubsub.Bus
	cfg pubsub.Config
}

func NewConsumer(bus *pubsub.Bus) *Consumer {
	c := new(Consumer)
	c.bus = bus
	// 🔥 configure first
	c.bus.SubConfig(c.configure)
	c.bus.SubOutbound(c.consume)
	return c
}

func (c *Consumer) configure(cfg pubsub.Config) {
	c.cfg = cfg
}

func (c *Consumer) consume(bytes []byte) {
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "yellow",
		EBCDIC: true,
		Title:  "Outbound",
	}
	c.bus.PubDump(dmp)
}

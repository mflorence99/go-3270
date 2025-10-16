package inbound

import (
	"go3270/emulator/pubsub"
)

type Producer struct {
	bus *pubsub.Bus
	cfg pubsub.Config
}

func NewProducer(bus *pubsub.Bus) *Producer {
	p := new(Producer)
	p.bus = bus
	// 👇 subscriptions
	p.bus.SubConfig(p.configure)
	return p
}

func (p *Producer) configure(cfg pubsub.Config) {
	p.cfg = cfg
}

func (p *Producer) Produce() {
	bytes := make([]byte, 0)
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "palegreen",
		EBCDIC: true,
		Title:  "Inbound",
	}
	p.bus.PubDump(dmp)
	p.bus.PubInbound(bytes)
}

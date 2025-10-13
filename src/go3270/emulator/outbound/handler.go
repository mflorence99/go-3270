package outbound

import (
	"go3270/emulator/pubsub"
)

type Handler struct {
	bus *pubsub.Bus
}

func NewHandler(bus *pubsub.Bus) *Handler {
	o := new(Handler)
	o.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	o.bus.SubClose(o.close)
	o.bus.SubOutbound(o.handle)
	return o
}

func (o *Handler) close() {}

func (o *Handler) handle(bytes []byte) {
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "yellow",
		EBCDIC: true,
		Title:  "Outbound",
	}
	o.bus.PubDump(dmp)
}

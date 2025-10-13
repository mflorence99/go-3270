package inbound

import (
	"go3270/emulator/pubsub"
)

type Handler struct {
	bus *pubsub.Bus
}

func NewHandler(bus *pubsub.Bus) *Handler {
	i := new(Handler)
	i.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	i.bus.SubClose(i.close)
	i.bus.SubInbound(i.handle)
	return i
}

func (i *Handler) close() {}

func (i *Handler) handle(bytes []byte) {
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "palegreen",
		EBCDIC: true,
		Title:  "Inbound",
	}
	i.bus.PubDump(dmp)
	i.bus.PubInbound(bytes)
}

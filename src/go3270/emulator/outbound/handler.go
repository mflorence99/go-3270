package outbound

import (
	"go3270/emulator/pubsub"
)

type Handler struct {
	bus *pubsub.Bus
}

func NewHandler(bus *pubsub.Bus) *Handler {
	k := new(Handler)
	k.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	k.bus.Subscribe(pubsub.CLOSE, k.close)
	k.bus.Subscribe(pubsub.OUTBOUND, k.handle)
	return k
}

func (k *Handler) close() {}

func (k *Handler) handle(bytes []byte) {
}

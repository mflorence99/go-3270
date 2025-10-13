package keyboard

import (
	"fmt"
	"go3270/emulator/pubsub"
)

// 🟧 Keyboard handler

type Handler struct {
	bus *pubsub.Bus
}

func NewHandler(bus *pubsub.Bus) *Handler {
	k := new(Handler)
	k.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	k.bus.SubClose(k.close)
	k.bus.SubKeystroke(k.handle)
	return k
}

func (k *Handler) close() {}

func (k *Handler) handle(key pubsub.Keystroke) {
	println(fmt.Sprintf("⌨️ %s", key))
}

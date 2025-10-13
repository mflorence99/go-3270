package keyboard

import (
	"fmt"
	"go3270/emulator/bus"
)

// 🟧 Keyboard handler

type Handler struct {
	bus *bus.Bus
}

func NewHandler(bus *bus.Bus) *Handler {
	k := new(Handler)
	k.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	k.bus.Subscribe("close", k.close)
	k.bus.Subscribe("keystroke", k.handle)
	return k
}

func (k *Handler) close() {}

func (k *Handler) handle(key Keystroke) {
	println(fmt.Sprintf("⌨️ %s", key))
}

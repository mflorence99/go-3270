package keyboard

import (
	"fmt"
	"go3270/emulator/pubsub"
)

type Consumer struct {
	bus *pubsub.Bus
}

func NewConsumer(bus *pubsub.Bus) *Consumer {
	k := new(Consumer)
	k.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	k.bus.SubClose(k.close)
	k.bus.SubKeystroke(k.consume)
	return k
}

func (k *Consumer) close() {}

func (k *Consumer) consume(key pubsub.Keystroke) {
	println(fmt.Sprintf("⌨️ %s", key))
}

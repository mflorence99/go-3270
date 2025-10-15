package keyboard

import (
	"fmt"
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
	c.bus.SubKeystroke(c.consume)
	return c
}

func (c *Consumer) configure(cfg pubsub.Config) {
	c.cfg = cfg
}

func (c *Consumer) consume(key pubsub.Keystroke) {
	println(fmt.Sprintf("⌨️ %s", key))
}

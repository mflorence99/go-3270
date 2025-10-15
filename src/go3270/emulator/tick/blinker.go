package tick

import (
	"go3270/emulator/pubsub"
)

type Blinker struct {
	bus *pubsub.Bus
	cfg pubsub.Config
}

func NewBlinker(bus *pubsub.Bus) *Blinker {
	c := new(Blinker)
	c.bus = bus
	// 🔥 configure first
	c.bus.SubConfig(c.configure)
	c.bus.SubTick(c.blink)
	return c
}

func (c *Blinker) blink(counter int) {}

func (c *Blinker) configure(cfg pubsub.Config) {
	c.cfg = cfg
}

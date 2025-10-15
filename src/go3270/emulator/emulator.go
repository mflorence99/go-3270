package emulator

import (
	"go3270/emulator/buffer"
	"go3270/emulator/glyph"
	"go3270/emulator/inbound"
	"go3270/emulator/keyboard"
	"go3270/emulator/outbound"
	"go3270/emulator/pubsub"
	"go3270/emulator/screen"
	"go3270/emulator/state"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	bus *pubsub.Bus
	buf *buffer.Buffer
	gc  *glyph.Cache
	in  *inbound.Producer
	key *keyboard.Consumer
	out *outbound.Consumer
	scr *screen.Screen
	st  *state.State
}

func NewEmulator(bus *pubsub.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	e.bus.SubClose(e.close)
	e.bus.SubConfig(e.configure)
	return e
}

func (e *Emulator) close() {}

func (e *Emulator) configure(cfg pubsub.Config) {
	e.buf = buffer.NewBuffer(cfg)
	e.gc = glyph.NewCache(cfg)
	e.in = inbound.NewProducer(e.bus)
	e.key = keyboard.NewConsumer(e.bus)
	e.out = outbound.NewConsumer(e.bus)
	e.scr = screen.NewScreen(cfg)
	e.st = state.NewState(e.bus)
}

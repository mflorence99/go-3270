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
	"go3270/emulator/tick"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	blnkr *tick.Blinker
	bus   *pubsub.Bus
	buf   *buffer.Buffer
	gc    *glyph.Cache
	in    *inbound.Producer
	key   *keyboard.Consumer
	out   *outbound.Consumer
	scr   *screen.Screen
	st    *state.State
}

func NewEmulator(bus *pubsub.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	e.blnkr = tick.NewBlinker(e.bus)
	e.buf = buffer.NewBuffer(e.bus)
	e.gc = glyph.NewCache(e.bus)
	e.in = inbound.NewProducer(e.bus)
	e.key = keyboard.NewConsumer(e.bus)
	e.out = outbound.NewConsumer(e.bus)
	e.scr = screen.NewScreen(e.bus)
	e.st = state.NewState(e.bus)
	// 🔥 configure first
	e.bus.SubConfig(e.configure)
	e.bus.SubClose(e.close)
	return e
}

func (e *Emulator) close() {
	println("🐞 Emulator closed")
}

func (e *Emulator) configure(cfg pubsub.Config) {}

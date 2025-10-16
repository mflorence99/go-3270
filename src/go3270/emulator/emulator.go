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
	blnkr *screen.Blinker
	bus   *pubsub.Bus
	buf   *buffer.Buffer
	gc    *glyph.Cache
	in    *inbound.Producer
	key   *keyboard.Keyboard
	out   *outbound.Consumer
	scr   *screen.Screen
	st    *state.State
}

func NewEmulator(bus *pubsub.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	// 👇 core components
	e.buf = buffer.NewBuffer(e.bus)
	e.gc = glyph.NewCache(e.bus)
	e.st = state.NewState(e.bus)
	// 👇 rendering components
	e.scr = screen.NewScreen(e.bus, e.buf, e.gc)
	e.blnkr = screen.NewBlinker(e.bus, e.buf, e.gc, e.scr, e.st)
	// 👇 i/o components
	e.key = keyboard.NewKeyboard(e.bus, e.buf, e.st)
	e.in = inbound.NewProducer(e.bus)
	e.out = outbound.NewConsumer(e.bus, e.buf, e.st)
	// 🔥 configure first
	e.bus.SubConfig(e.configure)
	e.bus.SubClose(e.close)
	return e
}

func (e *Emulator) close() {
	println("🐞 Emulator closed")
}

func (e *Emulator) configure(cfg pubsub.Config) {}

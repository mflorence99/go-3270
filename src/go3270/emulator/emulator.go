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
	cfg pubsub.Config
	gc  *glyph.Cache
	in  *inbound.Producer
	key *keyboard.Keyboard
	out *outbound.Consumer
	scr *screen.Screen
	st  *state.State
}

func NewEmulator(bus *pubsub.Bus, cfg pubsub.Config) *Emulator {
	e := new(Emulator)
	e.bus = bus
	e.cfg = cfg
	// 👇 core components; need these FIRST
	e.buf = buffer.NewBuffer(e.bus)
	e.gc = glyph.NewCache(e.bus)
	e.st = state.NewState(e.bus)
	// 👇 rendering components
	e.scr = screen.NewScreen(e.bus, e.buf, e.gc, e.st)
	// 👇 i/o components
	e.key = keyboard.NewKeyboard(e.bus, e.buf, e.st)
	e.in = inbound.NewProducer(e.bus)
	e.out = outbound.NewConsumer(e.bus, e.buf, e.st)
	// 👇 subscriptions
	e.bus.SubClose(e.close)
	// 👇 now configure all components
	e.bus.PubConfig(cfg)
	e.bus.PubReset()
	return e
}

func (e *Emulator) close() {
	println("🐞 Emulator closed")
}

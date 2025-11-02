package emulator

import (
	"go3270/emulator/buffer"
	"go3270/emulator/debug"
	"go3270/emulator/glyph"
	"go3270/emulator/inbound"
	"go3270/emulator/keyboard"
	"go3270/emulator/outbound"
	"go3270/emulator/pubsub"
	"go3270/emulator/screen"
	"go3270/emulator/screenshots"
	"go3270/emulator/state"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	bus   *pubsub.Bus
	buf   *buffer.Buffer
	cells *buffer.Cells
	cfg   pubsub.Config
	flds  *buffer.Flds
	gc    *glyph.Cache
	log   *debug.Logger
	in    *inbound.Producer
	key   *keyboard.Keyboard
	out   *outbound.Consumer
	scr   *screen.Screen
	st    *state.State
}

// 🟦 Constructor

func NewEmulator(bus *pubsub.Bus, cfg pubsub.Config) *Emulator {
	e := new(Emulator)
	e.bus = bus
	e.cfg = cfg
	// 🔥 absolutely #1 as FIFO for pubsub topics
	e.buf = buffer.NewBuffer(e.bus)
	// 👇 core components; need these FIRST
	e.gc = glyph.NewCache(e.bus)
	e.st = state.NewState(e.bus)
	e.cells = buffer.NewCells(e.bus, e.buf)
	e.flds = buffer.NewFlds(e.bus, e.buf)
	e.log = debug.NewLogger(e.bus, e.buf, e.flds, e.st)
	// 👇 rendering components
	e.scr = screen.NewScreen(e.bus, e.buf, e.gc, e.st)
	// 👇 i/o components
	e.key = keyboard.NewKeyboard(e.bus, e.buf, e.flds, e.st)
	e.in = inbound.NewProducer(e.bus, e.buf, e.flds, e.st)
	e.out = outbound.NewConsumer(e.bus, e.buf, e.cells, e.flds, e.st)
	// 👇 subscriptions
	e.bus.SubClose(e.close)
	// 👇 now configure all components
	e.bus.PubConfig(cfg)
	e.bus.PubReset()
	// 👇 if debugging, show screenshot
	if e.cfg.Screenshot != "" {
		e.bus.PubOutbound(screenshots.Index[e.cfg.Screenshot])
	}
	return e
}

// TODO 🔥 placeholder, just in case we need it

func (e *Emulator) close() {}

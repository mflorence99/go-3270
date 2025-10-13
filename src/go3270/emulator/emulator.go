package emulator

import (
	"go3270/emulator/buffer"
	"go3270/emulator/inbound"
	"go3270/emulator/keyboard"
	"go3270/emulator/outbound"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	bus *pubsub.Bus
	buf *buffer.Buffer
	in  *inbound.Handler
	key *keyboard.Handler
	out *outbound.Handler
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
	e.buf = buffer.NewBuffer(cfg.Rows * cfg.Cols)
	e.in = inbound.NewHandler(e.bus)
	e.key = keyboard.NewHandler(e.bus)
	e.out = outbound.NewHandler(e.bus)
	e.st = state.NewState(e.bus)
}

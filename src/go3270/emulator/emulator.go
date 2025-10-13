package emulator

import (
	"go3270/emulator/buffer"
	"go3270/emulator/keyboard"
	"go3270/emulator/pubsub"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	bus *pubsub.Bus
	buf *buffer.Buffer
	key *keyboard.Handler
}

func NewEmulator(bus *pubsub.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	e.bus.Subscribe(pubsub.CLOSE, e.close)
	e.bus.Subscribe(pubsub.CONFIG, e.configure)
	return e
}

func (e *Emulator) close() {}

func (e *Emulator) configure(cfg *Config) {
	e.buf = buffer.NewBuffer(cfg.Rows * cfg.Cols)
	e.key = keyboard.NewHandler(e.bus)
}

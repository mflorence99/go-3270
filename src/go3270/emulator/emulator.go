package emulator

import (
	"go3270/emulator/buffer"
	"go3270/emulator/bus"
	"go3270/emulator/keyboard"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	bus *bus.Bus
	buf *buffer.Buffer
	key *keyboard.Handler
}

func NewEmulator(bus *bus.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	// 🔥 must subscribe BEFORE we create any children
	e.bus.Subscribe("close", e.close)
	e.bus.Subscribe("config", e.configure)
	return e
}

func (e *Emulator) close() {}

func (e *Emulator) configure(cfg *Config) {
	e.buf = buffer.NewBuffer(cfg.Rows * cfg.Cols)
	e.key = keyboard.NewHandler(e.bus)
}

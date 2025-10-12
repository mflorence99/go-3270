package emulator

import (
	"go3270/emulator/bus"
)

type Emulator struct {
	bus *bus.Bus
}

func NewEmulator(bus *bus.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	return e
}

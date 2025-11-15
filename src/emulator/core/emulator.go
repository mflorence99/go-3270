package core

import (
	_ "embed"

	"emulator/types"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	Bus *Bus
	Cfg *types.Config
}

// 🟦 Constructor

func NewEmulator(bus *Bus, cfg *types.Config) *Emulator {
	e := new(Emulator)
	e.Bus = bus
	e.Cfg = cfg
	return e
}

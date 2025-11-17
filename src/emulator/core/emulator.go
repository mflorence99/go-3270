package core

import (
	_ "embed"

	"emulator/samples"
	"emulator/types"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	Buf    *Buffer
	Bus    *Bus
	Cache  *Cache
	Cells  *Cells
	Cfg    *types.Config
	Flds   *Flds
	Screen *Screen
	State  *State
}

// 🟦 Constructor

func NewEmulator(bus *Bus, cfg *types.Config) *Emulator {
	e := new(Emulator)
	e.Bus = bus
	e.Cfg = cfg
	// 🔥 preserve order of components for pubsub!
	e.Buf = NewBuffer(e)
	e.Cache = NewCache(e)
	e.Cells = NewCells(e)
	e.Flds = NewFlds(e)
	e.Screen = NewScreen(e)
	e.State = NewState(e)
	// 👇 subscriptions
	e.Bus.SubClose(e.close)
	// 👇 now initialize all components
	e.Bus.PubInit()
	// 👇 if debugging, show screenshot
	if e.Cfg.Testpage != "" {
		e.Bus.PubOutbound(samples.Index[e.Cfg.Testpage])
	}
	return e
}

// TODO 🔥 placeholder, just in case we need it

func (e *Emulator) close() {}

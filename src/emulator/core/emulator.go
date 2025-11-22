package core

import (
	_ "embed"

	"emulator/samples"
	"emulator/types"
)

// 🟧 3270 emulator itself, in pure go test-able code

type Emulator struct {
	Buf   *Buffer
	Bus   *Bus
	Cells *Cells
	Cfg   *types.Config
	Flds  *Flds
	GC    *Cache
	Kbd   *Keyboard
	In    *Producer
	Log   *Logger
	Out   *Consumer
	Scr   *Screen
	State *State
}

// 🟦 Constructor

func NewEmulator(bus *Bus, cfg *types.Config) *Emulator {
	e := new(Emulator)
	e.Bus = bus
	e.Cfg = cfg
	// 🔥 preserve order of components for pubsub!
	e.Buf = NewBuffer(e)
	e.Cells = NewCells(e)
	e.Flds = NewFlds(e)
	e.GC = NewCache(e)
	e.In = NewProducer(e)
	e.Kbd = NewKeyboard(e)
	e.Log = NewLogger(e)
	e.Out = NewConsumer(e)
	e.Scr = NewScreen(e)
	e.State = NewState(e)
	// 👇 subscriptions
	e.Bus.SubClose(e.close)
	return e
}

// TODO 🔥 placeholder, just in case we need it

func (e *Emulator) close() {}

// 🔥 caller initializes when ready

func (e *Emulator) Init() *Emulator {
	e.Bus.PubInit()
	// 👇 if debugging, show screenshot
	if e.Cfg.Testpage != "" {
		e.Bus.PubOutbound(samples.Index[e.Cfg.Testpage])
	}
	// 👇 useful for chaining directly to ctor
	return e
}

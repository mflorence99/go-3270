package screen

import (
	"go3270/emulator/buffer"
	"go3270/emulator/glyph"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"
)

type Blinker struct {
	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
	gc  *glyph.Cache
	scr *Screen
	st  *state.State
}

func NewBlinker(bus *pubsub.Bus, buf *buffer.Buffer, gc *glyph.Cache, scr *Screen, st *state.State) *Blinker {
	b := new(Blinker)
	b.bus = bus
	b.buf = buf
	b.gc = gc
	b.scr = scr
	b.st = st
	// 👇 subscriptions
	b.bus.SubConfig(b.configure)
	b.bus.SubTick(b.blink)
	return b
}

func (b *Blinker) blink(counter int) {
	blinkOn := counter%2 == 0
	// 👇 find all the blinbers
	blinkers := utils.NewStack[int](1)
	for addr := 0; addr < b.buf.Len(); addr++ {
		cell, _ := b.buf.Peek(addr)
		if !cell.FldStart && cell.Attrs.Blink {
			blinkers.Push(addr)
		}
	}
	// 👇 include the cursor if we have the focus
	if !b.st.Stat.Locked || !blinkOn {
		blinkers.Push(b.st.Stat.CursorAt)
	}
	// 👇 now we can render
	b.bus.PubBlink(blinkers, blinkOn)
}

func (b *Blinker) configure(cfg pubsub.Config) {
	b.cfg = cfg
}

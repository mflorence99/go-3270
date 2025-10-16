package screen

import (
	"go3270/emulator/buffer"
	"go3270/emulator/glyph"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"

	"github.com/fogleman/gg"
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
	// 🔥 configure first
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
	// 👇 iterate over all the blinkers
	dc := gg.NewContextForRGBA(b.cfg.RGBA)
	for !blinkers.Empty() {
		addr, _ := blinkers.Pop()
		// 👇 gather related data
		box := b.scr.CPs[addr]
		cell, _ := b.buf.Peek(addr)
		attrs := cell.Attrs
		// 👇 different color if highlighted
		color := utils.Ternary(attrs.Color == 0, b.cfg.Color, b.cfg.CLUT[attrs.Color])
		ix := utils.Ternary(attrs.Highlight, 1, 0)
		// 👇 the cache will find us the glyph iself
		g := glyph.Glyph{
			Char:       cell.Char,
			Color:      color[ix],
			Highlight:  attrs.Highlight,
			Reverse:    blinkOn,
			Underscore: attrs.Underscore,
		}
		img := b.gc.ImageFor(g, box)
		dc.DrawImage(img, int(box.X), int(box.Y))
	}
	// 👇 now we can render
	b.bus.PubRender()
}

func (b *Blinker) configure(cfg pubsub.Config) {
	b.cfg = cfg
}

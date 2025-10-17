package screen

import (
	"go3270/emulator/buffer"
	"go3270/emulator/glyph"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"

	"github.com/fogleman/gg"
)

type Screen struct {
	CPs []pubsub.Box

	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
	gc  *glyph.Cache
	st  *state.State
}

func NewScreen(bus *pubsub.Bus, buf *buffer.Buffer, gc *glyph.Cache, st *state.State) *Screen {
	s := new(Screen)
	s.bus = bus
	s.buf = buf
	s.gc = gc
	s.st = st
	// 👇 subscriptions
	s.bus.SubTick(s.blink)
	s.bus.SubConfig(s.configure)
	s.bus.SubRender(s.render)
	s.bus.SubReset(s.reset)
	return s
}

func (s *Screen) blink(counter int) {
	blinkOn := counter%2 == 0
	// 👇 find all the blinbers
	blinkers := utils.NewStack[int](1)
	for addr := 0; addr < s.buf.Len(); addr++ {
		cell, _ := s.buf.Peek(addr)
		if !cell.FldStart && cell.Attrs.Blink {
			blinkers.Push(addr)
		}
	}
	// 👇 include the cursor if we have the focus
	if !s.st.Stat.Locked || !blinkOn {
		blinkers.Push(s.st.Stat.CursorAt)
	}
	// 👇 now we can render
	s.renderImpl(blinkers, true, blinkOn)
}

func (s *Screen) configure(cfg pubsub.Config) {
	s.cfg = cfg
	s.CPs = make([]pubsub.Box, cfg.Cols*cfg.Rows)
	for ix := range s.CPs {
		row := int(ix / cfg.Cols)
		col := ix % cfg.Cols
		s.CPs[ix] = pubsub.NewBox(row, col, cfg)
	}
}

func (s *Screen) render(addrs *utils.Stack[int]) {
	s.renderImpl(addrs, false, false)
}

func (s *Screen) renderImpl(addrs *utils.Stack[int], doBlink bool, blinkOn bool) {
	// defer utils.ElapsedTime(time.Now())
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	// 👇 iterate over all changed cells
	for !addrs.Empty() {
		addr, _ := addrs.Pop()
		// 👇 gather related data
		box := s.CPs[addr]
		cell, _ := s.buf.Peek(addr)
		attrs := cell.Attrs
		invisible := cell.Char == 0x00 || cell.FldStart || attrs.Hidden
		// 👇 different color if highlighted
		color := utils.Ternary(attrs.Color == 0, s.cfg.Color, s.cfg.CLUT[attrs.Color])
		ix := utils.Ternary(attrs.Highlight, 1, 0)
		// 🔥 != here is the Go idion for XOR
		reverse := utils.Ternary(doBlink, attrs.Reverse != blinkOn, attrs.Reverse != (addr == s.st.Stat.CursorAt))
		// 👇 the cache will find us the glyph iself
		g := glyph.Glyph{
			Char:       utils.Ternary(invisible, ' ', cell.Char),
			Color:      color[ix],
			Highlight:  attrs.Highlight,
			Reverse:    reverse,
			Underscore: attrs.Underscore,
		}
		img := s.gc.ImageFor(g, box)
		dc.DrawImage(img, int(box.X), int(box.Y))
	}
}

func (s *Screen) reset() {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	dc.SetHexColor(s.cfg.BgColor)
	dc.Clear()
}

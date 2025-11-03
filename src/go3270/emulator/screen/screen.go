package screen

import (
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
	"go3270/emulator/glyph"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"

	"github.com/fogleman/gg"
)

// 🟧 Model the screen onto which the buffer is rendered (eventually an HTML <canvas> in the Typescript UI)

type Screen struct {
	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
	gc  *glyph.Cache
	st  *state.State

	clean  bool
	cps    []glyph.Box
	glyphs []glyph.Glyph
}

// 🟦 Constructor

func NewScreen(bus *pubsub.Bus, buf *buffer.Buffer, gc *glyph.Cache, st *state.State) *Screen {
	s := new(Screen)
	s.buf = buf
	s.bus = bus
	s.gc = gc
	s.st = st
	// 👇 subscriptions
	s.bus.SubTick(s.blink)
	s.bus.SubConfig(s.configure)
	s.bus.SubRender(s.render)
	// 🔥 curry the general function with the right parameters
	s.bus.SubRenderDeltas(func(deltas *utils.Stack[int]) { s.renderDeltas(deltas, false, false) })
	s.bus.SubReset(s.reset)
	return s
}

func (s *Screen) configure(cfg pubsub.Config) {
	s.cfg = cfg
	// 👇 precompute the box for each cell
	s.cps = make([]glyph.Box, s.cfg.Cols*s.cfg.Rows)
	for ix := range s.cps {
		row := int(ix / cfg.Cols)
		col := ix % cfg.Cols
		s.cps[ix] = glyph.NewBox(row, col, cfg)
	}
	// 👇 optimization remembers which glyph is already drawn in each cell
	s.glyphs = make([]glyph.Glyph, s.cfg.Cols*s.cfg.Rows)
}

func (s *Screen) reset() {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	dc.SetHexColor(s.cfg.BgColor)
	dc.Clear()
	s.glyphs = make([]glyph.Glyph, s.cfg.Cols*s.cfg.Rows)
	s.clean = true
}

// 🟦 Rendering functions

func (s *Screen) blink(counter int) {
	blinkOn := counter%2 == 0
	// 👇 find all the blinkers
	blinkers := utils.NewStack[int](1)
	for addr := 0; addr < s.buf.Len(); addr++ {
		cell, _ := s.buf.Peek(addr)
		if !cell.FldStart && cell.Attrs.Blink {
			blinkers.Push(addr)
		}
	}
	// 👇 include the cursor if we have the focus
	if !s.st.Status.Locked {
		blinkers.Push(s.st.Status.CursorAt)
	}
	// 👇 now we can render
	s.renderDeltas(blinkers, true, blinkOn)
}

func (s *Screen) render() {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	// 👇 iterate over all requested cells
	for addr := 0; addr < s.buf.Len(); addr++ {
		s.renderImpl(dc, addr, false, false)
	}
	s.clean = false
}

func (s *Screen) renderDeltas(addrs *utils.Stack[int], doBlink bool, blinkOn bool) {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	// 👇 iterate over all requested cells
	for !addrs.Empty() {
		addr, _ := addrs.Pop()
		s.renderImpl(dc, addr, doBlink, blinkOn)
	}
	s.clean = false
}

func (s *Screen) renderImpl(dc *gg.Context, addr int, doBlink bool, blinkOn bool) {
	// 👇 gather related data
	box := s.cps[addr]
	cell, _ := s.buf.Peek(addr)
	sf, _ := s.buf.Peek(cell.FldAddr)
	invisible := cell.Char == 0x00 || cell.FldStart || cell.Attrs.Hidden
	// 👇 ignore color if monochrome
	ix := utils.Ternary(cell.Attrs.Color == 0 || s.cfg.Monochrome, 0xF4, cell.Attrs.Color)
	color := s.cfg.CLUT[ix]
	highlight := cell.Attrs.Highlight
	// 🔥 outlined field can't be reverse or underscvore
	reverse := cell.Attrs.Reverse && sf.Attrs.Outline != 0x00
	underscore := cell.Attrs.Underscore && sf.Attrs.Outline != 0x00
	// 🔥 != is the Go idiom for XOR
	reverse = utils.Ternary(doBlink, reverse != blinkOn, reverse != (addr == s.st.Status.CursorAt))
	char := utils.Ternary(invisible, ' ', cell.Char)
	// 🔥 optimization: if the screen is clean and the char blank, skip
	if !s.clean || char > ' ' || reverse || cell.Attrs.Outline != 0x00 {
		// 👇 the cache will find us the glyph iself
		g := glyph.Glyph{
			Char:       char,
			Color:      color,
			Highlight:  highlight,
			Reverse:    reverse,
			Underscore: underscore,
		}
		// 🔥 outline is always at field level
		g.Outline.Bottom = (sf.Attrs.Outline & consts.OUTLINE_BOTTOM) != 0
		g.Outline.Right = ((sf.Attrs.Outline & consts.OUTLINE_RIGHT) != 0) && cell.FldEnd
		g.Outline.Top = (sf.Attrs.Outline & consts.OUTLINE_TOP) != 0
		g.Outline.Left = ((sf.Attrs.Outline & consts.OUTLINE_LEFT) != 0) && cell.FldStart
		// 👇 if the glyph is already at this address, no need to redraw it
		if g != s.glyphs[addr] {
			img := s.gc.ImageFor(g, box)
			dc.DrawImage(img, int(box.X), int(box.Y))
			s.glyphs[addr] = g
		}
	}
}

package screen

import (
	"go3270/emulator/buffer"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/types"
	"go3270/emulator/utils"

	"github.com/fogleman/gg"
)

// 🟧 Model the screen onto which the buffer is rendered
//    (eventually an HTML <canvas> in the Typescript UI)

type Screen struct {
	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg types.Config
	gc  *Cache
	st  *state.State

	clean  bool
	cps    []Box
	glyphs []Glyph
}

// 🟦 Constructor

func NewScreen(bus *pubsub.Bus, buf *buffer.Buffer, gc *Cache, st *state.State) *Screen {
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

func (s *Screen) configure(cfg types.Config) {
	s.cfg = cfg
	// 👇 precompute the box for each cell
	s.cps = make([]Box, s.cfg.Cols*s.cfg.Rows)
	for ix := range s.cps {
		row := int(ix / cfg.Cols)
		col := ix % cfg.Cols
		s.cps[ix] = NewBox(row, col, cfg)
	}
	// 👇 optimization remembers which glyph is already drawn in each cell
	s.glyphs = make([]Glyph, s.cfg.Cols*s.cfg.Rows)
}

func (s *Screen) reset() {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	dc.SetHexColor(s.cfg.BgColor)
	dc.Clear()
	s.glyphs = make([]Glyph, s.cfg.Cols*s.cfg.Rows)
	s.clean = true
}

// 🟦 Rendering functions

func (s *Screen) blink(counter int) {
	blinkOn := counter%2 == 1
	// 👇 find all the blinkers
	blinkers := utils.NewStack[int](1)
	for addr := 0; addr < s.buf.Len(); addr++ {
		cell := s.buf.MustPeek(addr)
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
	// 👇 iterate over all cells
	for addr := 0; addr < s.buf.Len(); addr++ {
		s.renderImpl(dc, addr, false, false)
	}
	s.clean = false
}

func (s *Screen) renderDeltas(addrs *utils.Stack[int], doBlink bool, blinkOn bool) {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	// 👇 iterate over all requested cells
	for !addrs.Empty() {
		addr, ok := addrs.Pop()
		if ok {
			s.renderImpl(dc, addr, doBlink, blinkOn)
		}
	}
	s.clean = false
}

func (s *Screen) renderImpl(dc *gg.Context, addr int, doBlink bool, blinkOn bool) {
	// 👇 gather related data
	box := s.cps[addr]
	cell := s.buf.MustPeek(addr)
	a := cell.Attrs
	// 👇 ignore color if monochrome
	var ix types.Color
	if s.cfg.Monochrome {
		ix = 0xf4
	} else if a.Color == 0x00 {
		switch {
		case !a.Protected && (a.Highlight || a.Hidden):
			a.Color = 0xF2
		case !a.Protected && !a.Highlight:
			a.Color = 0xF4
		case a.Protected && (a.Highlight || a.Hidden):
			a.Color = 0xF7
		case a.Protected && !a.Highlight:
			a.Color = 0xF1
		}
	} else {
		ix = a.Color
	}
	color := s.cfg.CLUT[ix]
	// 🔥 outlined field can't be reverse or underscore
	sf := s.buf.MustPeek(cell.FldAddr)
	fa := sf.Attrs
	outline := fa.Outline != 0x00
	reverse := a.Reverse && !outline
	underscore := a.Underscore && !outline && !cell.FldStart
	// 🔥 != is the Go idiom for XOR
	reverse = utils.Ternary(doBlink, reverse != blinkOn, reverse != (addr == s.st.Status.CursorAt))
	invisible := cell.Char == 0x00 || cell.FldStart || a.Hidden
	char := utils.Ternary(invisible, ' ', cell.Char)
	// 🔥 optimization: if the screen is clean and the char blank, skip
	if !s.clean || char > ' ' || outline || reverse || underscore {
		// 👇 the cache will find us the glyph iself
		g := Glyph{
			Char:       char,
			Color:      color,
			Highlight:  a.Highlight,
			Reverse:    reverse,
			Underscore: underscore,
			LCID:       a.LCID,
		}
		// 🔥 outline is always at field level
		if outline {
			g.Outline.Bottom = (fa.Outline & types.OUTLINE_BOTTOM) != 0
			g.Outline.Right = ((fa.Outline & types.OUTLINE_RIGHT) != 0) && cell.FldEnd
			g.Outline.Top = (fa.Outline & types.OUTLINE_TOP) != 0
			g.Outline.Left = ((fa.Outline & types.OUTLINE_LEFT) != 0) && cell.FldStart
		}
		// 👇 if the glyph is already at this address, no need to redraw it
		if g != s.glyphs[addr] {
			img := s.gc.ImageFor(g, box)
			dc.DrawImage(img, int(box.X), int(box.Y))
			s.glyphs[addr] = g
		}
	}
}

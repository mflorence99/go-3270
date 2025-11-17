package core

import (
	"emulator/types"
	"emulator/utils"

	"github.com/fogleman/gg"
)

// 🟧 Model the screen onto which the buffer is rendered
//    (eventually an HTML <canvas> in the Typescript UI)

type Screen struct {
	clean  bool
	cps    []Box
	glyphs []Glyph

	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewScreen(emu *Emulator) *Screen {
	s := new(Screen)
	s.emu = emu
	// 👇 subscriptions
	s.emu.Bus.SubInit(s.init)
	s.emu.Bus.SubRender(s.render)
	s.emu.Bus.SubTick(s.blink)
	// 🔥 curry the general function with the right parameters
	s.emu.Bus.SubRenderDeltas(func(deltas *utils.Stack[uint]) {
		s.renderDeltas(deltas, false, false)
	})
	s.emu.Bus.SubReset(s.reset)
	return s
}

func (s *Screen) init() {
	// 👇 precompute the box for each cell
	s.cps = make([]Box, s.emu.Cfg.Cols*s.emu.Cfg.Rows)
	for ix := range s.cps {
		row := uint(ix) / s.emu.Cfg.Cols
		col := uint(ix) % s.emu.Cfg.Cols
		s.cps[ix] = NewBox(row, col, s.emu.Cfg)
	}
	// 👇 optimization remembers which glyph is already drawn in each cell
	s.glyphs = make([]Glyph, s.emu.Cfg.Cols*s.emu.Cfg.Rows)
}

func (s *Screen) reset() {
	dc := gg.NewContextForRGBA(s.emu.Cfg.RGBA)
	dc.SetHexColor(s.emu.Cfg.BgColor)
	dc.Clear()
	s.glyphs = make([]Glyph, s.emu.Cfg.Cols*s.emu.Cfg.Rows)
	s.clean = true
}

// 🟦 Rendering functions

func (s *Screen) blink(counter int) {
	blinkOn := counter%2 == 1
	// 👇 find all the blinkers
	blinkers := utils.NewStack[uint](1)
	for addr := uint(0); addr < s.emu.Buf.Len(); addr++ {
		cell := s.emu.Buf.MustPeek(addr)
		if !cell.IsFldStart() && cell.Attrs.Blink {
			blinkers.Push(addr)
		}
	}
	// 👇 include the cursor if we have the focus
	if !s.emu.State.Status.Locked {
		blinkers.Push(s.emu.State.Status.CursorAt)
	}
	// 👇 now we can render
	s.renderDeltas(blinkers, true, blinkOn)
}

func (s *Screen) render() {
	dc := gg.NewContextForRGBA(s.emu.Cfg.RGBA)
	// 👇 iterate over all cells
	for addr := uint(0); addr < s.emu.Buf.Len(); addr++ {
		s.renderImpl(dc, addr, false, false)
	}
	s.clean = false
}

func (s *Screen) renderDeltas(addrs *utils.Stack[uint], doBlink bool, blinkOn bool) {
	dc := gg.NewContextForRGBA(s.emu.Cfg.RGBA)
	// 👇 iterate over all requested cells
	for !addrs.Empty() {
		addr, ok := addrs.Pop()
		if ok {
			s.renderImpl(dc, addr, doBlink, blinkOn)
		}
	}
	s.clean = false
}

func (s *Screen) renderImpl(dc *gg.Context, addr uint, doBlink bool, blinkOn bool) {
	// 👇 gather related data
	box := s.cps[addr]
	cell := s.emu.Buf.MustPeek(addr)
	a := cell.Attrs
	color := s.emu.Cfg.ColorOf(a)
	// 🔥 outlined field can't be reverse or underscore and must be on field
	sf, ok := cell.GetFldStart()
	var outline types.Outline
	if ok {
		outline = sf.Attrs.Outline
	}
	reverse := a.Reverse && outline == 0x00
	underscore := a.Underscore && outline == 0x00 && !cell.IsFldStart()
	// 🔥 != is the Go idiom for XOR
	reverse = utils.Ternary(doBlink, reverse != blinkOn, reverse != (addr == s.emu.State.Status.CursorAt))
	invisible := cell.Char == 0x00 || cell.IsFldStart() || a.Hidden
	char := utils.Ternary(invisible, ' ', cell.Char)
	// 🔥 optimization: if the screen is clean and the char blank, skip
	if !s.clean || char > ' ' || outline != 0x00 || reverse || underscore {
		// 👇 the cache will find us the glyph iself
		g := Glyph{
			Char:       char,
			Color:      color,
			Highlight:  a.Highlight,
			Reverse:    reverse,
			Underscore: underscore,
			LCID:       a.LCID,
		}
		// 👇 outline surrounds the entire field
		if outline != 0b00000000 {
			fld, _ := cell.FindFld()
			g.Outline = Outline{
				Bottom: (outline & types.OUTLINE_BOTTOM) != 0,
				Right:  ((outline & types.OUTLINE_RIGHT) != 0) && cell == fld.Cells[len(fld.Cells)-1],
				Top:    (outline & types.OUTLINE_TOP) != 0,
				Left:   ((outline & types.OUTLINE_LEFT) != 0) && cell == fld.Cells[0],
			}
		}
		// 👇 if the glyph is already at this address, no need to redraw it
		if g != s.glyphs[addr] {
			img := s.emu.Cache.ImageFor(g, box)
			dc.DrawImage(img, int(box.X), int(box.Y))
			s.glyphs[addr] = g
		}
	}
}

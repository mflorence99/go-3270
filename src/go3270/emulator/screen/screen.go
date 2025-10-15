package screen

import (
	"go3270/emulator/buffer"
	"go3270/emulator/glyph"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"time"

	"github.com/fogleman/gg"
)

type Screen struct {
	CPs []pubsub.Box

	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
	gc  *glyph.Cache
}

func NewScreen(bus *pubsub.Bus, buf *buffer.Buffer, gc *glyph.Cache) *Screen {
	s := new(Screen)
	s.bus = bus
	s.buf = buf
	s.gc = gc
	// 🔥 configure first
	s.bus.SubConfig(s.configure)
	s.bus.SubRender(s.render)
	s.bus.SubReset(s.reset)
	return s
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

func (s *Screen) render() {
	defer utils.ElapsedTime(time.Now())
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	// 👇 iterate over all changed cells
	for !s.buf.Dirty.Empty() {
		addr, _ := s.buf.Dirty.Pop()
		// 👇 gather related data
		box := s.CPs[addr]
		cell, _ := s.buf.Peek(addr)
		attrs := cell.Attrs
		// 👇 different color if highlighted
		color := s.cfg.CLUT[attrs.Color]
		ix := utils.Ternary(attrs.Highlight, 1, 0)
		// 👇 the cache will find us the glyph iself
		g := glyph.Glyph{
			Char:       cell.Char,
			Color:      color[ix],
			Highlight:  attrs.Highlight,
			Reverse:    attrs.Reverse,
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

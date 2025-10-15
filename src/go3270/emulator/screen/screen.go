package screen

import (
	"go3270/emulator/pubsub"

	"github.com/fogleman/gg"
)

type Screen struct {
	CPs []Box

	bus *pubsub.Bus
	cfg pubsub.Config
}

func NewScreen(bus *pubsub.Bus) *Screen {
	s := new(Screen)
	s.bus = bus
	// 🔥 configure first
	s.bus.SubConfig(s.configure)
	s.bus.SubReset(s.reset)
	return s
}

func (s *Screen) configure(cfg pubsub.Config) {
	s.cfg = cfg
	s.CPs = make([]Box, cfg.Cols*cfg.Rows)
	for ix := range s.CPs {
		row := int(ix / cfg.Cols)
		col := ix % cfg.Cols
		s.CPs[ix] = NewBox(row, col, cfg)
	}
}

func (s *Screen) reset() {
	dc := gg.NewContextForRGBA(s.cfg.RGBA)
	dc.SetHexColor(s.cfg.BgColor)
	dc.Clear()
}

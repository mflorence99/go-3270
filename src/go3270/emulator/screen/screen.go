package screen

import (
	"go3270/emulator/pubsub"
)

type Screen struct {
	CPs []Box
}

func NewScreen(cfg pubsub.Config) *Screen {
	s := new(Screen)
	s.CPs = make([]Box, cfg.Cols*cfg.Rows)
	for ix := range s.CPs {
		row := int(ix / cfg.Cols)
		col := ix % cfg.Cols
		s.CPs[ix] = NewBox(row, col, cfg)
	}
	return s
}

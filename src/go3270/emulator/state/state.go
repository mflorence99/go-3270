package state

import (
	"fmt"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
)

type State struct {
	bus  *pubsub.Bus
	cfg  pubsub.Config
	stat pubsub.Status
}

func NewState(bus *pubsub.Bus) *State {
	s := new(State)
	s.bus = bus
	// 🔥 configure first
	s.bus.SubConfig(s.configure)
	s.bus.SubReset(s.reset)
	return s
}

func (s *State) configure(cfg pubsub.Config) {
	s.cfg = cfg
}

func (s *State) reset() {
	s.Patch(Patch{
		Alarm:     utils.BoolPtr(false),
		CursorAt:  utils.IntPtr(0),
		Error:     utils.BoolPtr(false),
		Locked:    utils.BoolPtr(false),
		Message:   utils.StringPtr(""),
		Numeric:   utils.BoolPtr(false),
		Protected: utils.BoolPtr(false),
		Waiting:   utils.BoolPtr(false),
	})
}

func (s *State) Patch(p Patch) {
	if p.Alarm != nil {
		s.stat.Alarm = *p.Alarm
	}
	if p.CursorAt != nil {
		s.stat.CursorAt = *p.CursorAt
	}
	if p.Error != nil {
		s.stat.Error = *p.Error
	}
	if p.Locked != nil {
		s.stat.Locked = *p.Locked
	}
	if p.Message != nil {
		s.stat.Message = *p.Message
	}
	if p.Numeric != nil {
		s.stat.Numeric = *p.Numeric
	}
	if p.Protected != nil {
		s.stat.Protected = *p.Protected
	}
	if p.Waiting != nil {
		s.stat.Waiting = *p.Waiting
	}
	s.bus.PubStatus(s.stat)
	println(fmt.Sprintf("⚙️ %s", s.stat))
	// 👇 make sure to reset alarm
	s.stat.Alarm = false
}

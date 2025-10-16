package state

import (
	"fmt"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
)

type State struct {
	Stat *pubsub.Status

	bus *pubsub.Bus
	cfg pubsub.Config
}

func NewState(bus *pubsub.Bus) *State {
	s := new(State)
	s.bus = bus
	// 👇 subscriptions
	s.bus.SubConfig(s.configure)
	s.bus.SubReset(s.reset)
	return s
}

func (s *State) configure(cfg pubsub.Config) {
	s.cfg = cfg
	s.Stat = &pubsub.Status{}
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
		s.Stat.Alarm = *p.Alarm
	}
	if p.CursorAt != nil {
		s.Stat.CursorAt = *p.CursorAt
	}
	if p.Error != nil {
		s.Stat.Error = *p.Error
	}
	if p.Locked != nil {
		s.Stat.Locked = *p.Locked
	}
	if p.Message != nil {
		s.Stat.Message = *p.Message
	}
	if p.Numeric != nil {
		s.Stat.Numeric = *p.Numeric
	}
	if p.Protected != nil {
		s.Stat.Protected = *p.Protected
	}
	if p.Waiting != nil {
		s.Stat.Waiting = *p.Waiting
	}
	s.bus.PubStatus(s.Stat)
	println(fmt.Sprintf("⚙️ %s", s.Stat))
	// 👇 make sure to reset alarm
	s.Stat.Alarm = false
}

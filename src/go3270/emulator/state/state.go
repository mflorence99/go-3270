package state

import (
	"fmt"
	"go3270/emulator/pubsub"
)

type State struct {
	bus  *pubsub.Bus
	stat pubsub.Status
}

func NewState(bus *pubsub.Bus) *State {
	s := new(State)
	s.bus = bus
	return s
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

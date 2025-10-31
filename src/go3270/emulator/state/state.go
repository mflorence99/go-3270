package state

import (
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"go3270/emulator/wcc"
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
	s.bus.SubInbound(s.lock)
	s.bus.SubOutbound(s.unlock)
	s.bus.SubReset(s.reset)
	s.bus.SubWCC(s.wcc)
	return s
}

func (s *State) configure(cfg pubsub.Config) {
	s.cfg = cfg
	s.Stat = &pubsub.Status{}
}

func (s *State) lock(_ []byte, _ bool) {
	s.Patch(Patch{
		Locked:  utils.BoolPtr(true),
		Waiting: utils.BoolPtr(true),
	})
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

func (s *State) unlock(_ []byte) {
	s.Patch(Patch{
		Locked:  utils.BoolPtr(false),
		Waiting: utils.BoolPtr(false),
	})
}

func (s *State) wcc(wcc wcc.WCC) {
	// 👇 honor WCC instructions
	s.Patch(Patch{
		Alarm:  utils.BoolPtr(wcc.Alarm),
		Locked: utils.BoolPtr(!wcc.Unlock),
	})
}

// 🟦 Public methods

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
	// 👇 make sure to reset alarm
	s.Stat.Alarm = false
}

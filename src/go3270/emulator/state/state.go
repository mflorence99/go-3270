package state

import (
	"go3270/emulator/pubsub"
	"go3270/emulator/types"
	"go3270/emulator/utils"
)

// 🟧 3270 status (as shared with the Typescript UI)

type State struct {
	Status *types.Status

	bus *pubsub.Bus
	cfg types.Config
}

// 🟦 Constructor

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

func (s *State) configure(cfg types.Config) {
	s.cfg = cfg
	s.Status = &types.Status{}
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

// 🟦 Functions to dispatch actions depending on state

func (s *State) lock(_ []byte, _ pubsub.InboundHints) {
	s.Patch(Patch{
		Locked:  utils.BoolPtr(true),
		Waiting: utils.BoolPtr(true),
	})
}

func (s *State) unlock(_ []byte) {
	s.Patch(Patch{
		Locked:  utils.BoolPtr(false),
		Waiting: utils.BoolPtr(false),
	})
}

func (s *State) wcc(wcc types.WCC) {
	// 👇 honor WCC instructions
	s.Patch(Patch{
		Alarm:  utils.BoolPtr(wcc.Alarm),
		Locked: utils.BoolPtr(!wcc.Unlock),
	})
}

// 🟦 Public functions

func (s *State) Patch(p Patch) {
	if p.Alarm != nil {
		s.Status.Alarm = *p.Alarm
	}
	if p.CursorAt != nil {
		s.Status.CursorAt = *p.CursorAt
	}
	if p.Error != nil {
		s.Status.Error = *p.Error
	}
	if p.Locked != nil {
		s.Status.Locked = *p.Locked
	}
	if p.Message != nil {
		s.Status.Message = *p.Message
	}
	if p.Numeric != nil {
		s.Status.Numeric = *p.Numeric
	}
	if p.Protected != nil {
		s.Status.Protected = *p.Protected
	}
	if p.Waiting != nil {
		s.Status.Waiting = *p.Waiting
	}
	s.bus.PubStatus(s.Status)
	// 👇 make sure to reset alarm
	s.Status.Alarm = false
}

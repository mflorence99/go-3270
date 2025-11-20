package core

import (
	"emulator/types"
	"emulator/utils"
)

// 🟧 3270 status (as shared with the Typescript UI)

type State struct {
	Status *types.Status

	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewState(emu *Emulator) *State {
	s := new(State)
	s.emu = emu
	// 👇 subscriptions
	s.emu.Bus.SubInit(s.init)
	s.emu.Bus.SubInbound(s.lock)
	s.emu.Bus.SubOutbound(s.unlock)
	s.emu.Bus.SubReset(s.reset)
	s.emu.Bus.SubWCC(s.wcc)
	return s
}

func (s *State) init() {
	s.Status = &types.Status{}
}

func (s *State) reset() {
	s.Patch(types.Patch{
		Alarm:     utils.BoolPtr(false),
		CursorAt:  utils.UintPtr(0),
		Error:     utils.BoolPtr(false),
		Insert:    utils.BoolPtr(false),
		Locked:    utils.BoolPtr(false),
		Message:   utils.StringPtr(""),
		Numeric:   utils.BoolPtr(false),
		Protected: utils.BoolPtr(false),
		Waiting:   utils.BoolPtr(false),
	})
}

// 🟦 Functions to dispatch actions depending on state

func (s *State) lock(_ []byte, _ InboundHints) {
	s.Patch(types.Patch{
		Locked:  utils.BoolPtr(true),
		Waiting: utils.BoolPtr(true),
	})
}

func (s *State) unlock(_ []byte) {
	s.Patch(types.Patch{
		Locked:  utils.BoolPtr(false),
		Waiting: utils.BoolPtr(false),
	})
}

func (s *State) wcc(wcc types.WCC) {
	// 👇 honor WCC instructions
	s.Patch(types.Patch{
		Alarm:  utils.BoolPtr(wcc.Alarm),
		Locked: utils.BoolPtr(!wcc.Unlock),
	})
}

// 🟦 Public functions

func (s *State) Patch(p types.Patch) {
	if p.Alarm != nil {
		s.Status.Alarm = *p.Alarm
	}
	if p.CursorAt != nil {
		s.Status.CursorAt = *p.CursorAt
	}
	if p.Error != nil {
		s.Status.Error = *p.Error
	}
	if p.Insert != nil {
		s.Status.Insert = *p.Insert
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
	s.emu.Bus.PubStatus(s.Status)
	// 👇 make sure to reset alarm
	s.Status.Alarm = false
}

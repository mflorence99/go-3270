package core

import (
	"emulator/types"
	"emulator/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	emu := MockEmulator().Initialize()
	t.Run("smoke test on empty Status", func(t *testing.T) {
		assert.False(t, emu.State.Status.Alarm)
		assert.Equal(t, uint(0), emu.State.Status.CursorAt)
		assert.False(t, emu.State.Status.Error)
		assert.False(t, emu.State.Status.Insert)
		assert.Empty(t, emu.State.Status.Message)
		assert.False(t, emu.State.Status.Numeric)
		assert.False(t, emu.State.Status.Protected)
		assert.False(t, emu.State.Status.Waiting)
		assert.False(t, emu.State.Status.Alarm)
	})
}

func TestStatePatch(t *testing.T) {
	emu := MockEmulator().Initialize()
	emu.State.Patch(types.Patch{
		Alarm:     utils.BoolPtr(true),
		CursorAt:  utils.UintPtr(100),
		Error:     utils.BoolPtr(true),
		Insert:    utils.BoolPtr(true),
		Locked:    utils.BoolPtr(true),
		Message:   utils.StringPtr("help!"),
		Numeric:   utils.BoolPtr(true),
		Protected: utils.BoolPtr(true),
		Waiting:   utils.BoolPtr(true),
	})
	t.Run("test status after patching", func(t *testing.T) {
		// 🔥 Alarm is reset after patch
		// assert.True(t, emu.State.Status.Alarm)
		assert.Equal(t, uint(100), emu.State.Status.CursorAt)
		assert.True(t, emu.State.Status.Error)
		assert.True(t, emu.State.Status.Insert)
		assert.True(t, emu.State.Status.Locked)
		assert.Equal(t, "help!", emu.State.Status.Message)
		assert.True(t, emu.State.Status.Numeric)
		assert.True(t, emu.State.Status.Protected)
		assert.True(t, emu.State.Status.Waiting)
	})
}

func TestStateWCC(t *testing.T) {
	emu := MockEmulator().Initialize()
	var actual types.WCC
	emu.Bus.SubWCChar(func(wcc types.WCC) {
		actual = wcc
	})
	expected := types.WCC{
		Alarm:    true,
		Reset:    true,
		ResetMDT: true,
		Unlock:   true,
	}
	emu.Bus.PubWCChar(expected)
	assert.Equal(t, expected, actual, "WCC recorded correctly")
}

func TestStateLock(t *testing.T) {
	emu := MockEmulator().Initialize()
	var locked, unlocked bool

	emu.Bus.SubInbound(func(_ []byte, _ PubInboundHints) {
		locked = true
		unlocked = false
	})
	emu.Bus.PubInbound(nil, PubInboundHints{})
	t.Run("status is locked after inbound", func(t *testing.T) {
		assert.True(t, locked)
		assert.False(t, unlocked)
	})

	emu.Bus.SubOutbound(func(_ []byte) {
		locked = false
		unlocked = true
	})
	emu.Bus.PubOutbound(nil)
	t.Run("status is unlocked after outbound", func(t *testing.T) {
		assert.False(t, locked)
		assert.True(t, unlocked)
	})
}

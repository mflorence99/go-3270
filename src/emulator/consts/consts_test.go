package consts_test

import (
	"emulator/consts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AIDFor(t *testing.T) {
	assert.True(t, consts.AIDFor(0x6C) == "PA1")
	assert.True(t, consts.AIDFor(0x6D) == "CLEAR")
	assert.True(t, consts.AIDFor(0xFF) == "")
}

func Test_AIDFOf(t *testing.T) {
	assert.True(t, consts.AIDOf("Enter", false, false, false) == consts.ENTER)
	assert.True(t, consts.AIDOf("Escape", false, false, false) == consts.CLEAR)
	assert.True(t, consts.AIDOf("F1", false, false, false) == consts.PF1)
	assert.True(t, consts.AIDOf("F3", false, false, false) == consts.PF3)
	assert.True(t, consts.AIDOf("F1", false, false, true) == consts.PF13)
	assert.True(t, consts.AIDOf("F3", false, false, true) == consts.PF15)
	assert.True(t, consts.AIDOf("F1", true, false, false) == consts.PA1)
	assert.True(t, consts.AIDOf("F3", true, false, false) == consts.PA3)
	assert.True(t, consts.AIDOf("F4", true, false, false) == 0)
	assert.True(t, consts.AIDOf("Backspace", true, false, false) == 0)
}

func Test_PAx(t *testing.T) {
	assert.True(t, consts.PAx(0x6C))
	assert.False(t, consts.PAx(0x00))
	assert.False(t, consts.PAx(0x88))
}

func Test_PFx(t *testing.T) {
	assert.True(t, consts.PFx(0x4C))
	assert.True(t, consts.PFx(0xC5))
	assert.False(t, consts.PFx(0x88))
}

func Test_CommandFor(t *testing.T) {
	assert.True(t, consts.CommandFor(0x6F) == "EAU")
	assert.True(t, consts.CommandFor(0xF5) == "EW")
	assert.True(t, consts.CommandFor(0xFF) == "")
}

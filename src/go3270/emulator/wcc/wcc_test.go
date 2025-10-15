package wcc_test

import (
	"go3270/emulator/wcc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWCC_New(t *testing.T) {
	var w = wcc.NewWCC(0b00000100)
	assert.True(t, w.Alarm)
	w = wcc.NewWCC(0b01000000)
	assert.True(t, w.Reset)
	w = wcc.NewWCC(0b00000001)
	assert.True(t, w.ResetMDT)
	w = wcc.NewWCC(0b00000010)
	assert.True(t, w.Unlock)
}

func TestWCC_Byte(t *testing.T) {
	var w = wcc.NewWCC(0b11111111)
	assert.True(t, w.Byte() == 0b01000111)
}

func TestWCC_String(t *testing.T) {
	var w = wcc.NewWCC(0b11111111)
	assert.True(t, w.String() == "WCC=[ ALARM RESET -MDT UNLOCK ]")
	w = wcc.NewWCC(0b00000000)
	assert.True(t, w.String() == "WCC=[ ]")
}

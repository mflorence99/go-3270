package wcc_test

import (
	"emulator/wcc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWCC_New(t *testing.T) {
	var ch = wcc.New(0b00000100)
	assert.True(t, ch.Alarm())
	ch = wcc.New(0b01000000)
	assert.True(t, ch.Reset())
	ch = wcc.New(0b00000001)
	assert.True(t, ch.ResetMDT())
	ch = wcc.New(0b00000010)
	assert.True(t, ch.Unlock())
}

func TestWCC_Byte(t *testing.T) {
	var ch = wcc.New(0b11111111)
	assert.True(t, ch.Byte() == 0b01000111)
}

func TestWCC_String(t *testing.T) {
	var ch = wcc.New(0b11111111)
	assert.True(t, ch.String() == "WCC=[ ALARM RESET -MDT UNLOCK ]")
	ch = wcc.New(0b00000000)
	assert.True(t, ch.String() == "WCC=[ ]")
}

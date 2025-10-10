package wcc_test

import (
	"emulator/wcc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWCC_New(t *testing.T) {
	var char = wcc.New(0b00000100)
	assert.True(t, char.Alarm())
	char = wcc.New(0b01000000)
	assert.True(t, char.Reset())
	char = wcc.New(0b00000001)
	assert.True(t, char.ResetMDT())
	char = wcc.New(0b00000010)
	assert.True(t, char.Unlock())
}

func TestWCC_Byte(t *testing.T) {
	var char = wcc.New(0b11111111)
	assert.True(t, char.Byte() == 0b01000111)
}

func TestWCC_String(t *testing.T) {
	var char = wcc.New(0b11111111)
	assert.True(t, char.String() == "WCC=[ ALARM RESET -MDT UNLOCK ]")
	char = wcc.New(0b00000000)
	assert.True(t, char.String() == "WCC=[ ]")
}

package device_test

import (
	"emulator/device"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWCC_New(t *testing.T) {
	var wcc = device.NewWCC(0b00000100)
	assert.True(t, wcc.DoAlarm())
	wcc = device.NewWCC(0b01000000)
	assert.True(t, wcc.DoReset())
	wcc = device.NewWCC(0b00000001)
	assert.True(t, wcc.DoResetMDT())
	wcc = device.NewWCC(0b00000010)
	assert.True(t, wcc.DoUnlock())
}

func TestWCC_ToByte(t *testing.T) {
	var wcc = device.NewWCC(0b11111111)
	assert.True(t, wcc.ToByte() == 0b01000111)
}

func TestWCC_ToString(t *testing.T) {
	var wcc = device.NewWCC(0b11111111)
	assert.True(t, wcc.ToString() == "WCC=[ ALARM RESET -MDT UNLOCK ]")
	wcc = device.NewWCC(0b00000000)
	assert.True(t, wcc.ToString() == "WCC=[ ]")
}

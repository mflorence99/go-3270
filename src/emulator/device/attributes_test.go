package device_test

import (
	"emulator/device"
	"testing"

	"github.com/stretchr/testify/assert"
)

var attributes = []uint8{0xC0, 0b00110001, 0x41, 0xF1, 0x41, 0xF2, 0x41, 0xF4, 0x42, 0xF4}

func TestAttributes_New(t *testing.T) {
	var attrs = device.NewAttributes(attributes)
	assert.True(t, attrs.GetColor("") == "#88DD88")
	assert.True(t, attrs.IsBlink())
	assert.True(t, attrs.IsModified())
	assert.True(t, attrs.IsNumeric())
	assert.True(t, attrs.IsProtected())
	assert.True(t, attrs.IsReverse())
	assert.True(t, attrs.IsUnderscore())

}

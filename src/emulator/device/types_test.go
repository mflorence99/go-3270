package device_test

import (
	"emulator/device"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes_AID(t *testing.T) {
	assert.True(t, device.AID[0x4C] == "PF24")
	assert.True(t, device.AIDLookup["PF8"] == 0xF8)
}

func TestTypes_Command(t *testing.T) {
	assert.True(t, device.Command[0xF1] == "W")
	assert.True(t, device.CommandLookup["EW"] == 0xF5)
}

func TestTypes_Highlight(t *testing.T) {
	assert.True(t, device.Highlight[0xF2] == "REVERSE")
	assert.True(t, device.HighlightLookup["UNDERSCORE"] == 0xF4)
}

func TestTypes_Op(t *testing.T) {
	assert.True(t, device.Op[0x02] == "Q")
	assert.True(t, device.OpLookup["RM"] == 0xF6)
}

func TestTypes_Order(t *testing.T) {
	assert.True(t, device.Order[0x1D] == "SF")
	assert.True(t, device.OrderLookup["EUA"] == 0x12)
}

func TestTypes_TypeCode(t *testing.T) {
	assert.True(t, device.TypeCode[0xC0] == "BASIC")
	assert.True(t, device.TypeCodeLookup["COLOR"] == 0x42)
}

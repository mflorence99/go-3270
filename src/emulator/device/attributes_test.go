package device_test

import (
	"emulator/device"
	"testing"

	"github.com/stretchr/testify/assert"
)

var attributes = []byte{0xC0, 0b00110001, 0x41, 0xF1, 0x41, 0xF2, 0x41, 0xF4, 0x42, 0xF4}

func TestAttributes_NewAttributes(t *testing.T) {
	attrs := device.NewAttributes(attributes)
	assert.True(t, attrs.Blink())
	assert.True(t, attrs.Modified())
	assert.True(t, attrs.Numeric())
	assert.True(t, attrs.Protected())
	assert.True(t, attrs.Reverse())
	assert.True(t, attrs.Underscore())

}

func TestAttributes_ToString(t *testing.T) {
	attrs := device.NewAttributes(attributes)
	str := attrs.ToString()
	assert.True(t, str == "ATTR=[ BLINK MDT NUM PROT REV USCORE ]")
}

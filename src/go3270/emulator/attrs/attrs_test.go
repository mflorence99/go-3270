package attrs_test

import (
	"go3270/emulator/attrs"
	"testing"

	"github.com/stretchr/testify/assert"
)

var basic byte = 0b00111001

var extended = []byte{0xC0, 0b00111001, 0x41, 0xF1, 0x41, 0xF2, 0x41, 0xF4, 0x42, 0xF4}

func Test_NewBasic(t *testing.T) {
	a := attrs.NewBasic(basic)
	assert.False(t, a.Blink)
	assert.True(t, a.Highlight)
	assert.True(t, a.Modified)
	assert.True(t, a.Numeric)
	assert.True(t, a.Protected)
	assert.False(t, a.Reverse)
	assert.False(t, a.Underscore)
}

func Test_NewExtended(t *testing.T) {
	a := attrs.NewExtended(extended)
	assert.True(t, a.Blink)
	assert.True(t, a.Highlight)
	assert.True(t, a.Modified)
	assert.True(t, a.Numeric)
	assert.True(t, a.Protected)
	assert.True(t, a.Reverse)
	assert.True(t, a.Underscore)
}

func Test_Byte(t *testing.T) {
	a := attrs.NewBasic(basic)
	byte := a.Byte()
	assert.True(t, byte == 0b00111001)
}

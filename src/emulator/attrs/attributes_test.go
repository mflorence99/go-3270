package attrs_test

import (
	"emulator/attrs"
	"testing"

	"github.com/stretchr/testify/assert"
)

var attribute = []byte{0b00111001}

var attributes = []byte{0xC0, 0b00111001, 0x41, 0xF1, 0x41, 0xF2, 0x41, 0xF4, 0x42, 0xF4}

func Test_New(t *testing.T) {
	attrs := attrs.New(attribute)
	assert.False(t, attrs.Blink())
	assert.True(t, attrs.Highlight())
	assert.True(t, attrs.Modified())
	assert.True(t, attrs.Numeric())
	assert.True(t, attrs.Protected())
	assert.False(t, attrs.Reverse())
	assert.False(t, attrs.Underscore())
}

func Test_NewExtended(t *testing.T) {
	attrs := attrs.New(attributes)
	assert.True(t, attrs.Blink())
	assert.True(t, attrs.Highlight())
	assert.True(t, attrs.Modified())
	assert.True(t, attrs.Numeric())
	assert.True(t, attrs.Protected())
	assert.True(t, attrs.Reverse())
	assert.True(t, attrs.Underscore())
}

func Test_Byte(t *testing.T) {
	attrs := attrs.New(attribute)
	ch := attrs.Byte()
	assert.True(t, ch == 0b00111001)
}

func Test_String(t *testing.T) {
	attrs := attrs.New(attributes)
	str := attrs.String()
	assert.True(t, str == "ATTR=[ BLINK HILITE MDT NUM PROT REV USCORE ]")
}

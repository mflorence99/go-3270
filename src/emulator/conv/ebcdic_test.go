package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2A(t *testing.T) {
	assert.Equal(t, E2A(0x21), byte(' '), "everything below 0x40 is ASCII blank")
	assert.Equal(t, E2A(0x40), byte(' '), "EBCDIC 0x40 is ASCII blank")
	assert.Equal(t, E2A(0xf0), byte('0'), "EBCDIC 0xf0 is ASCII 0")
}

func TestE2As(t *testing.T) {
	hello := E2As(string([]byte{200, 133, 147, 147, 150}))
	assert.Equal(t, hello, "Hello", "convert EBCDIC string to ASCII string")
}

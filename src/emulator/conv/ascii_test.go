package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestA2E(t *testing.T) {
	assert.Equal(t, A2E(' '), byte(0x40), "ASCII blank is EBCDIC 0x40")
	assert.Equal(t, A2E('0'), byte(0xF0), "ASCII 0 is EBCDIC 0xF0")
}

func TestA2Es(t *testing.T) {
	hello := string([]byte{200, 133, 147, 147, 150})
	assert.Equal(t, A2Es("Hello"), hello, "convert ASCII string to EBCDIC")
}

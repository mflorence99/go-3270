package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddr2Bytes(t *testing.T) {
	addr := uint(100)
	bytes := Addr2Bytes(addr)
	assert.Equal(t, []byte{0xc1, 0xe4}, bytes, "convert # to 3270 address")
}

func TestBytes2Addr(t *testing.T) {
	bytes := []byte{0xc1, 0xe4}
	addr := Bytes2Addr(bytes)
	assert.Equal(t, uint(100), addr, "convert 3270 address to #")
}

package inbound_test

import (
	"emulator/stream/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stream = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}

func Test_Put(t *testing.T) {
	in := inbound.New()
	slice := in.Put(0x00)
	assert.True(t, len(slice) == 1)
	assert.True(t, slice[0] == 0x00)
	slice = in.PutSlice(stream)
	assert.True(t, len(slice) == 7)
	assert.True(t, slice[0] == 0x00)
	assert.True(t, slice[6] == 0x05)
}

func Test_Put16(t *testing.T) {
	in := inbound.New()
	slice := in.Put16(0x1234)
	assert.True(t, slice[0] == 0x12)
	assert.True(t, slice[1] == 0x34)
}

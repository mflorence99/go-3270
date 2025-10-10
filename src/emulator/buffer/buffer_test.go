package buffer_test

import (
	"emulator/buffer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	b := buffer.New(100)
	assert.True(t, b != nil)
	c := buffer.Cell{Char: 0x00}
	assert.True(t, c.Char == 0x00)
}

func Test_Seek(t *testing.T) {
	b := buffer.New(100)
	addr, ok := b.Seek(99)
	assert.True(t, addr == 99)
	assert.True(t, ok)
	addr, ok = b.Seek(100)
	assert.True(t, addr == 0)
	assert.False(t, ok)
}

package buffer_test

import (
	"emulator/buffer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Len() == 100)
	c := buffer.Cell{Char: 0x00}
	assert.True(t, c.Char == 0x00)
}

func Test_Len(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Len() == 100)
}

func Test_Peek(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, ok := b.Peek(50)
	assert.True(t, c == nil)
	assert.True(t, ok)
	c, ok = b.Peek(150)
	assert.True(t, c == nil)
	assert.False(t, ok)
}

func Test_Seek(t *testing.T) {
	b := buffer.NewBuffer(100)
	addr, ok := b.Seek(99)
	assert.True(t, addr == 99)
	assert.True(t, ok)
	addr, ok = b.Seek(100)
	assert.True(t, addr == -1)
	assert.False(t, ok)
}

func Test_Get(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, addr := b.Get()
	assert.True(t, c == nil)
	assert.True(t, addr == 0)
}

func Test_GetNext(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, addr := b.GetNext()
	assert.True(t, c == nil)
	assert.True(t, addr == 1)
}

func Test_GetPrev(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, addr := b.PrevGet()
	assert.True(t, c == nil)
	assert.True(t, addr == 99)
}

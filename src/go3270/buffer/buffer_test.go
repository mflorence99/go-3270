package buffer_test

import (
	"go3270/attrs"
	"go3270/buffer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Len() == 100)
	c := buffer.Cell{Char: 0x00}
	assert.True(t, c.Char == 0x00)
}

// 🟦 Housekeeping methods

func Test_Len(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Len() == 100)
}

func Test_Peek(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, ok := b.Peek(50)
	assert.True(t, c.Char == 0x00)
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

// 🟦 Get methods

func Test_Get(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, addr := b.Get()
	assert.True(t, c.Char == 0x00)
	assert.True(t, addr == 0)
}

func Test_GetNext(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, addr := b.GetNext()
	assert.True(t, c.Char == 0x00)
	assert.True(t, addr == 1)
}

func Test_GetPrev(t *testing.T) {
	b := buffer.NewBuffer(100)
	c, addr := b.PrevGet()
	assert.True(t, c.Char == 0x00)
	assert.True(t, addr == 99)
}

// 🟦 Set methods

func makeCell(num bool) *buffer.Cell {
	c := &buffer.Cell{
		Attrs:    &attrs.Attrs{Numeric: num},
		Char:     0x40,
		FldAddr:  0,
		FldStart: false,
	}
	return c
}

func Test_Set(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Changes.Len() == 0)
	c := makeCell(false)
	addr := b.Set(c)
	assert.True(t, b.Changes.Len() == 1)
	assert.True(t, addr == 0)
	c, _ = b.Get()
	assert.True(t, c.Char == 0x40)
}

func Test_SetAndNext(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Changes.Len() == 0)
	b.Seek(99)
	c := makeCell(false)
	addr := b.SetAndNext(c)
	assert.True(t, b.Changes.Len() == 1)
	assert.True(t, addr == 99)
	c, _ = b.Get()
	assert.True(t, c.Char == 0x00)
}

func Test_StartFld(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Changes.Len() == 0)
	b.Seek(50)
	addr := b.StartFld(&attrs.Attrs{})
	assert.True(t, b.Changes.Len() == 1)
	assert.True(t, addr == 50)
}

func Test_PrevAndSet(t *testing.T) {
	b := buffer.NewBuffer(100)
	assert.True(t, b.Changes.Len() == 0)
	c := makeCell(false)
	addr := b.PrevAndSet(c)
	assert.True(t, b.Changes.Len() == 1)
	assert.True(t, addr == 99)
	c, _ = b.Get()
	assert.True(t, c.Char == 0x40)
}

// 🟦 Keystroke methods

func Test_Keyin(t *testing.T) {
	b := buffer.NewBuffer(100)
	addr, ok := b.Keyin('x')
	assert.True(t, addr == -1)
	assert.False(t, ok)
	b.StartFld(&attrs.Attrs{})
	b.SetAndNext(makeCell(true))
	b.SetAndNext(makeCell(true))
	b.SetAndNext(makeCell(true))
	assert.True(t, b.Changes.Len() == 4)
	b.Seek(3)
	addr, ok = b.Keyin('x')
	assert.True(t, addr == -1)
	assert.False(t, ok)
	addr, ok = b.Keyin('0')
	assert.True(t, addr == 3)
	assert.True(t, ok)
	c, ok := b.Peek(0)
	assert.True(t, c.Attrs.Modified)
	assert.True(t, ok)
}

func Test_Backspace(t *testing.T) {
	b := buffer.NewBuffer(100)
	b.StartFld(&attrs.Attrs{})
	b.SetAndNext(makeCell(false))
	b.SetAndNext(makeCell(false))
	b.SetAndNext(makeCell(false))
	addr, ok := b.Backspace()
	assert.True(t, addr == 3)
	assert.True(t, ok)
}

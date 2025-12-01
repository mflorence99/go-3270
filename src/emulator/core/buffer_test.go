package core

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {
	emu := MockEmulator(12, 40).Initialize()

	assert.Equal(t, uint(0), emu.Buf.Addr(), "initial state")
	assert.Equal(t, uint(12*40), emu.Buf.Len(), "initial state")
	assert.Equal(t, types.FIELD_MODE, emu.Buf.Mode(), "initial state")
}

func TestBufferSetMode(t *testing.T) {
	emu := MockEmulator(12, 40).Initialize()

	assert.Equal(t, types.FIELD_MODE, emu.Buf.Mode(), "initial state")
	emu.Buf.SetMode(types.CHARACTER_MODE)
	assert.Equal(t, types.CHARACTER_MODE, emu.Buf.Mode(), "set mode to character mode")
}

func TestBufferPeek(t *testing.T) {
	var addr uint
	var cell *Cell
	var ok bool

	emu := MockEmulator(12, 40).Initialize()

	cell, ok = emu.Buf.Peek(100)
	assert.NotNil(t, cell, "Peek addr in range")
	assert.True(t, ok, "Peek addr in range")

	cell, ok = emu.Buf.Peek(1000)
	assert.Nil(t, cell, "Peek addr outside range")
	assert.False(t, ok, "Peek addr outside range")

	cell, addr = emu.Buf.WrappingPeek(1000)
	assert.NotNil(t, cell, "WrappingPeek addr in range")
	assert.Equal(t, uint(40), addr, "WrappingPeek addr in range")

	emu.Bus.SubPanic(func(msg string) {
		assert.Contains(t, msg, "out of range", "MustPeek addr outside range")
	})
	cell = emu.Buf.MustPeek(1000)
	assert.Nil(t, cell, "MustPeek addr outside range")
}

func TestBufferSeek(t *testing.T) {
	var addr uint
	var ok bool

	emu := MockEmulator(12, 40).Initialize()

	addr, ok = emu.Buf.Seek(100)
	assert.Equal(t, uint(100), addr, "Seek addr in range")
	assert.True(t, ok, "Seek addr in range")

	addr = emu.Buf.WrappingSeek(1000)
	assert.Equal(t, uint(40), addr, "WrappingSeek addr in range")

	emu.Bus.SubPanic(func(msg string) {
		assert.Contains(t, msg, "out of range", "MustSeek addr outside range")
	})
	addr = emu.Buf.MustSeek(1000)
	assert.Equal(t, uint(0), addr, "MustSeek addr outside range")
}

func TestBufferReplace(t *testing.T) {
	var cell *Cell
	var ok bool
	var repl = &Cell{}

	emu := MockEmulator(12, 40).Initialize()

	cell, ok = emu.Buf.Replace(repl, 100)
	assert.Equal(t, *repl, *cell, "Replace addr in range")
	assert.True(t, ok, "Replace addr in range")

	cell, ok = emu.Buf.Replace(repl, 1000)
	assert.Nil(t, cell, "Replace addr outside range")
	assert.False(t, ok, "Replace addr outside range")

	emu.Bus.SubPanic(func(msg string) {
		assert.Contains(t, msg, "out of range", "MustReplace addr outside range")
	})
	cell = emu.Buf.MustReplace(repl, 1000)
	assert.Nil(t, cell, "MustReplace addr outside range")
}

func TestBefferGet(t *testing.T) {
	var addr uint
	var cell *Cell

	emu := MockEmulator(12, 40).Initialize()
	emu.Buf.Seek(uint(12*40 - 1))

	cell, addr = emu.Buf.Get()
	assert.NotNil(t, cell, "Get cell at end of buffer")
	assert.Equal(t, uint(12*40-1), addr, "Get cell at end of buffer")

	cell, addr = emu.Buf.GetNext()
	assert.NotNil(t, cell, "Get cell after")
	assert.Equal(t, uint(0), addr, "Get cell after")

	cell, addr = emu.Buf.PrevGet()
	assert.NotNil(t, cell, "Get cell before")
	assert.Equal(t, uint(12*40-2), addr, "Get cell before")
}

func TestBufferSet(t *testing.T) {
	var addr uint
	var cell *Cell
	var repl = &Cell{}

	emu := MockEmulator(12, 40).Initialize()
	emu.Buf.Seek(uint(12*40 - 1))

	addr = emu.Buf.Set(repl)
	cell = emu.Buf.MustPeek(addr)
	assert.Equal(t, *repl, *cell, "Set cell correctly")

	emu.Buf.SetAndNext(repl)
	addr = emu.Buf.Addr()
	assert.Equal(t, uint(0), addr)

	emu.Buf.PrevAndSet(repl)
	addr = emu.Buf.Addr()
	assert.Equal(t, uint(12*40-1), addr)
}

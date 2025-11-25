package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOutbound(t *testing.T) {
	bus := NewBus()
	out := NewOutbound([]byte{}, bus)
	t.Run("smoke test on an empty stream", func(t *testing.T) {
		var char byte
		var short uint16
		var slice []byte
		var ok bool

		assert.Empty(t, out.Bytes())
		assert.True(t, out.HasEnough(0))
		assert.False(t, out.HasEnough(1))
		assert.False(t, out.HasNext())

		char, ok = out.Next()
		assert.Equal(t, char, byte(0))
		assert.False(t, ok)

		short, ok = out.Next16()
		assert.Equal(t, short, uint16(0))
		assert.False(t, ok)

		slice, ok = out.NextSlice(0)
		assert.Empty(t, slice)
		assert.True(t, ok)

		slice, ok = out.NextSlice(1)
		assert.Empty(t, slice)
		assert.False(t, ok)

		char, ok = out.Peek()
		assert.Equal(t, char, byte(0))
		assert.False(t, ok)

		short, ok = out.Peek16()
		assert.Equal(t, short, uint16(0))
		assert.False(t, ok)

		slice, ok = out.PeekSlice(0)
		assert.Empty(t, slice)
		assert.True(t, ok)

		slice, ok = out.PeekSlice(1)
		assert.Empty(t, slice)
		assert.False(t, ok)

		slice = out.Rest()
		assert.Empty(t, slice)

		ok = out.Skip(0)
		assert.True(t, ok)

		ok = out.Skip(1)
		assert.False(t, ok)
	})
}

func TestOutboundNext(t *testing.T) {
	bus := NewBus()
	out := NewOutbound([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, bus)
	t.Run("next functions with data", func(t *testing.T) {
		var char byte
		var short uint16
		var slice []byte
		var ok bool

		assert.NotEmpty(t, out.Bytes())
		assert.True(t, out.HasEnough(6))
		assert.False(t, out.HasEnough(7))
		assert.True(t, out.HasNext())

		char, ok = out.Next()
		assert.Equal(t, char, byte(0x00))
		assert.True(t, ok)

		short, ok = out.Next16()
		assert.Equal(t, short, uint16(0x0102))
		assert.True(t, ok)

		slice, ok = out.NextSlice(1)
		assert.Equal(t, slice, []byte{0x03})
		assert.True(t, ok)

		slice = out.Rest()
		assert.Equal(t, slice, []byte{0x04, 0x05})

		char, ok = out.Peek()
		assert.Equal(t, char, byte(0))
		assert.False(t, ok)
	})
}

func TestOutboundMustNext(t *testing.T) {
	bus := NewBus()

	numPanics := 0
	bus.SubPanic(func(msg string) {
		numPanics++
	})

	out := NewOutbound([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, bus)
	t.Run("next functions with data", func(t *testing.T) {

		_ = out.Skip(6)

		_ = out.MustNext()
		assert.Equal(t, 1, numPanics)

		_ = out.MustNext16()
		assert.Equal(t, 2, numPanics)

		_ = out.MustNextSlice(1)
		assert.Equal(t, 3, numPanics)
	})
}

func TestOutboundPeek(t *testing.T) {
	bus := NewBus()
	out := NewOutbound([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, bus)
	t.Run("next functions with data", func(t *testing.T) {
		var char byte
		var short uint16
		var slice []byte
		var ok bool

		char, ok = out.Peek()
		assert.Equal(t, char, byte(0x00))
		assert.True(t, ok)

		short, ok = out.Peek16()
		assert.Equal(t, short, uint16(0x0001))
		assert.True(t, ok)

		slice, ok = out.PeekSlice(3)
		assert.Equal(t, slice, []byte{0x00, 0x01, 0x02})
		assert.True(t, ok)

		slice, ok = out.PeekSliceUntil([]byte{0x04, 0x05})
		assert.Equal(t, slice, []byte{0x00, 0x01, 0x02, 0x03})
		assert.True(t, ok)

		_ = out.Skip(2)
		slice, ok = out.PeekSlice(3)
		assert.Equal(t, slice, []byte{0x02, 0x03, 0x04})
		assert.True(t, ok)
	})
}

package outbound_test

import (
	"emulator/stream/outbound"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stream = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}

func Test_HasEnough(t *testing.T) {
	out := outbound.New(&stream)
	assert.True(t, out.HasEnough(0))
	assert.True(t, out.HasEnough(5))
	assert.True(t, out.HasEnough(6))
	assert.False(t, out.HasEnough(7))
}

func Test_Next(t *testing.T) {
	out := outbound.New(&stream)
	for ix := 0; ix <= len(stream); ix++ {
		char, ok := out.Next()
		if ix < len(stream) {
			assert.True(t, char <= 6)
			assert.True(t, ok)
		} else {
			assert.True(t, char == 0)
			assert.True(t, !ok)
		}
	}
}

func Test_Next16(t *testing.T) {
	out := outbound.New(&stream)
	chars, ok := out.Next16()
	assert.True(t, chars == 0x01)
	assert.True(t, ok)
	chars, ok = out.Next16()
	assert.True(t, chars == 0x0203)
	assert.True(t, ok)
	chars, ok = out.Next16()
	assert.True(t, chars == 0x0405)
	assert.True(t, ok)
	chars, ok = out.Next16()
	assert.True(t, chars == 0)
	assert.True(t, !ok)
}

func Test_NextSlice(t *testing.T) {
	out := outbound.New(&stream)
	_, ok := out.Next()
	assert.True(t, ok)
	char, ok := out.Next()
	assert.True(t, char == 0x01)
	assert.True(t, ok)
	slice, ok := out.NextSlice(4)
	assert.True(t, slices.Equal(slice, []byte{0x02, 0x03, 0x04, 0x05}))
	assert.True(t, ok)
	slice, ok = out.NextSlice(4)
	assert.True(t, len(slice) == 0)
	assert.True(t, !ok)
}

func Test_NextSliceUntil(t *testing.T) {
	out := outbound.New(&stream)
	slice, ok := out.NextSliceUntil([]byte{0x02, 0x03})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01}))
	assert.True(t, ok)
	slice, ok = out.NextSliceUntil([]byte{0x03, 0x04})
	assert.True(t, slices.Equal(slice, []byte{0x02}))
	assert.True(t, ok)
	out.Skip(2)
	slice, ok = out.NextSliceUntil([]byte{0x06, 0x07})
	assert.True(t, slices.Equal(slice, []byte{0x05}))
	assert.True(t, !ok)
}

func Test_Peek(t *testing.T) {
	out := outbound.New(&stream)
	char, _ := out.Peek()
	assert.True(t, char == 0x00)
	slice, ok := out.NextSlice(6)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, ok)
	char, ok = out.Peek()
	assert.True(t, char == 0)
	assert.True(t, !ok)
}

func Test_PeekSlice(t *testing.T) {
	out := outbound.New(&stream)
	slice, ok := out.PeekSlice(6)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, ok)
	slice, ok = out.PeekSlice(2)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01}))
	assert.True(t, ok)
	slice, ok = out.PeekSlice(7)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, !ok)
}

func Test_PeekSliceUntil(t *testing.T) {
	out := outbound.New(&stream)
	slice, ok := out.PeekSliceUntil([]byte{0x02, 0x03})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01}))
	assert.True(t, ok)
	slice, ok = out.PeekSliceUntil([]byte{0x03, 0x04})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02}))
	assert.True(t, ok)
	slice, ok = out.PeekSliceUntil([]byte{0x06, 0x07})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, !ok)
}

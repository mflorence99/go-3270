package device_test

import (
	"emulator/device"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stream = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}

func TestOutboundDataStream_HasEnough(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	assert.True(t, out.HasEnough(0))
	assert.True(t, out.HasEnough(5))
	assert.True(t, out.HasEnough(6))
	assert.False(t, out.HasEnough(7))
}

func TestOutboundDataStream_Next(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	for ix := 0; ix <= len(stream); ix++ {
		u8, e := out.Next()
		if ix < len(stream) {
			assert.True(t, u8 <= 6)
			assert.True(t, e == nil)
		} else {
			assert.True(t, u8 == 0)
			assert.True(t, e != nil)
		}
	}
}

func TestOutboundDataStream_Next16(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	var u16 uint16
	var err error
	u16, err = out.Next16()
	assert.True(t, u16 == 0x01)
	assert.True(t, err == nil)
	u16, err = out.Next16()
	assert.True(t, u16 == 0x0203)
	assert.True(t, err == nil)
	u16, err = out.Next16()
	assert.True(t, u16 == 0x0405)
	assert.True(t, err == nil)
	u16, err = out.Next16()
	assert.True(t, u16 == 0)
	assert.True(t, err != nil)
}

func TestOutboundDataStream_NextSlice(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	var u8 byte
	var slice []byte
	var err error
	_, _ = out.Next()
	u8, _ = out.Next()
	assert.True(t, u8 == 0x01)
	slice, err = out.NextSlice(4)
	assert.True(t, slices.Equal(slice, []byte{0x02, 0x03, 0x04, 0x05}))
	assert.True(t, err == nil)
	slice, err = out.NextSlice(4)
	assert.True(t, len(slice) == 0)
	assert.True(t, err != nil)
}

func TestOutboundDataStream_NextSliceUntil(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	var slice []byte
	var err error
	slice, err = out.NextSliceUntil([]byte{0x02, 0x03})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01}))
	assert.True(t, err == nil)
	slice, err = out.NextSliceUntil([]byte{0x03, 0x04})
	assert.True(t, slices.Equal(slice, []byte{0x02}))
	assert.True(t, err == nil)
	out.Skip(2)
	slice, err = out.NextSliceUntil([]byte{0x06, 0x07})
	assert.True(t, slices.Equal(slice, []byte{0x05}))
	assert.True(t, err != nil)
}

func TestOutboundDataStream_Peek(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	var u8 byte
	var slice []byte
	var err error
	u8, _ = out.Peek()
	assert.True(t, u8 == 0x00)
	slice, err = out.NextSlice(6)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, err == nil)
	u8, err = out.Peek()
	assert.True(t, u8 == 0)
	assert.True(t, err != nil)
}

func TestOutboundDataStream_PeekSlice(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	var slice []byte
	var err error
	slice, err = out.PeekSlice(6)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, err == nil)
	slice, err = out.PeekSlice(2)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01}))
	assert.True(t, err == nil)
	slice, err = out.PeekSlice(7)
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, err != nil)
}

func TestOutboundDataStream_PeekSliceUntil(t *testing.T) {
	out := device.NewOutboundDataStream(&stream)
	var slice []byte
	var err error
	slice, err = out.PeekSliceUntil([]byte{0x02, 0x03})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01}))
	assert.True(t, err == nil)
	slice, err = out.PeekSliceUntil([]byte{0x03, 0x04})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02}))
	assert.True(t, err == nil)
	slice, err = out.PeekSliceUntil([]byte{0x06, 0x07})
	assert.True(t, slices.Equal(slice, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.True(t, err != nil)
}

func TestInboundDataStream_Put(t *testing.T) {
	in := device.NewInboundDataStream()
	var slice []byte
	slice = in.Put(0x00)
	assert.True(t, len(slice) == 1)
	assert.True(t, slice[0] == 0x00)
	slice = in.PutSlice(stream)
	assert.True(t, len(slice) == 7)
	assert.True(t, slice[0] == 0x00)
	assert.True(t, slice[6] == 0x05)
}

func TestInboundDataStream_Put16(t *testing.T) {
	in := device.NewInboundDataStream()
	slice := in.Put16(0x1234)
	assert.True(t, slice[0] == 0x12)
	assert.True(t, slice[1] == 0x34)
}

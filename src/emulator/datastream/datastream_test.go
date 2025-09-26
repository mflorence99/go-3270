package datastream_test

import (
	"emulator/datastream"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = []uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}

func TestOutboundDataStream_HasEnough(t *testing.T) {
	out := datastream.NewOutbound(&data)
	assert.True(t, out.HasEnough(0))
	assert.True(t, out.HasEnough(5))
	assert.True(t, out.HasEnough(6))
	assert.False(t, out.HasEnough(7))
}

func TestOutboundDataStream_Next(t *testing.T) {
	out := datastream.NewOutbound(&data)
	for ix := 0; ix <= len(data); ix++ {
		u8, e := out.Next()
		if ix < len(data) {
			assert.True(t, u8 <= 6)
			assert.True(t, e == nil)
		} else {
			assert.True(t, u8 == 0)
			assert.True(t, e != nil)
		}
	}
}

func TestOutboundDataStream_Next16(t *testing.T) {
	out := datastream.NewOutbound(&data)
	var u16 uint16
	var e any
	u16, e = out.Next16()
	assert.True(t, u16 == 0x01)
	assert.True(t, e == nil)
	u16, e = out.Next16()
	assert.True(t, u16 == 0x0203)
	assert.True(t, e == nil)
	u16, e = out.Next16()
	assert.True(t, u16 == 0x0405)
	assert.True(t, e == nil)
	u16, e = out.Next16()
	assert.True(t, u16 == 0)
	assert.True(t, e != nil)
}

func TestOutboundDataStream_Slice(t *testing.T) {
	out := datastream.NewOutbound(&data)
	var u8 uint8
	var slice []uint8
	var e any
	_, e = out.Next()
	u8, _ = out.Next()
	assert.True(t, u8 == 0x01)
	slice, e = out.NextSlice(4)
	assert.True(t, slice[0] == 0x02 && slice[3] == 0x05)
	assert.True(t, e == nil)
	slice, e = out.NextSlice(4)
	assert.True(t, slice == nil)
	assert.True(t, e != nil)
}

func TestOutboundDataStream_Peek(t *testing.T) {
	out := datastream.NewOutbound(&data)
	var u8 uint8
	var slice []uint8
	var e any
	u8, _ = out.Peek()
	assert.True(t, u8 == 0x00)
	slice, e = out.NextSlice(6)
	assert.True(t, slice[0] == 0x00 && slice[5] == 0x05)
	assert.True(t, e == nil)
	u8, e = out.Peek()
	assert.True(t, u8 == 0)
	assert.True(t, e != nil)
}

func TestOutboundDataStream_PeekSlice(t *testing.T) {
	out := datastream.NewOutbound(&data)
	var slice []uint8
	var e any
	slice, e = out.PeekSlice(6)
	assert.True(t, slice[0] == 0x00 && slice[5] == 0x05)
	assert.True(t, e == nil)
	slice, e = out.PeekSlice(2)
	assert.True(t, slice[0] == 0x00 && slice[1] == 0x01)
	assert.True(t, e == nil)
	slice, e = out.PeekSlice(7)
	assert.True(t, slice == nil)
	assert.True(t, e != nil)
}

func TestInboundDataStream_Put(t *testing.T) {
	in := datastream.NewInbound()
	var slice []uint8
	slice = in.Put(0x00)
	assert.True(t, len(slice) == 1)
	assert.True(t, slice[0] == 0x00)
	slice = in.PutSlice(data)
	assert.True(t, len(slice) == 7)
	assert.True(t, slice[0] == 0x00)
	assert.True(t, slice[6] == 0x05)
}

func TestInboundDataStream_Put16(t *testing.T) {
	in := datastream.NewInbound()
	slice := in.Put16(0x1234)
	assert.True(t, slice[0] == 0x12)
	assert.True(t, slice[1] == 0x34)
}

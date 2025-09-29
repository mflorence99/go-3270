package device

import (
	"errors"
	"slices"
)

// 🟧 Model outbound 3270 data as a stream
//    "Outbound" data flows from the application to the 3270 ie this code

type OutboundDataStream struct {
	bytes *[]uint8
	ix    int
}

func NewOutboundDataStream(bytes *[]uint8) *OutboundDataStream {
	stream := &OutboundDataStream{}
	stream.bytes = bytes
	stream.ix = 0
	return stream
}

func (out *OutboundDataStream) HasEnough(count int) bool {
	return out.ix+count-1 < len(*out.bytes)
}

func (out *OutboundDataStream) HasNext() bool {
	return out.ix < len(*out.bytes)
}

func (out *OutboundDataStream) Next() (uint8, error) {
	return out.nextImpl(false)
}

func (out *OutboundDataStream) Next16() (uint16, error) {
	hi, e1 := out.Next()
	if e1 != nil {
		return 0, e1
	}
	lo, e2 := out.Next()
	if e2 != nil {
		return 0, e2
	}
	return (uint16(hi) * 256) + uint16(lo), nil
}

func (out *OutboundDataStream) NextSlice(count int) ([]uint8, error) {
	return out.nextSliceImpl(count, false)
}

func (out *OutboundDataStream) NextSliceUntil(matches []uint8) ([]uint8, error) {
	return out.nextSliceUntilImpl(matches, false)
}

func (out *OutboundDataStream) Peek() (uint8, error) {
	return out.nextImpl(true)
}

func (out *OutboundDataStream) PeekSlice(count int) ([]uint8, error) {
	return out.nextSliceImpl(count, true)
}

func (out *OutboundDataStream) PeekSliceUntil(matches []uint8) ([]uint8, error) {
	return out.nextSliceUntilImpl(matches, true)
}

// 👇 Helpers

func (out *OutboundDataStream) nextImpl(peek bool) (uint8, error) {
	if out.HasNext() {
		u8 := (*out.bytes)[out.ix]
		if !peek {
			out.ix += 1
		}
		return u8, nil
	} else {
		return 0, errors.New("insufficient bytes in stream")
	}
}

func (out *OutboundDataStream) nextSliceImpl(count int, peek bool) ([]uint8, error) {
	if out.HasEnough(count) {
		end := out.ix + count
		slice := (*out.bytes)[out.ix:end]
		if !peek {
			out.ix = end
		}
		return slice, nil
	} else {
		return nil, errors.New("insufficient bytes in stream")
	}
}

func (out *OutboundDataStream) nextSliceUntilImpl(matches []uint8, peek bool) ([]uint8, error) {
	count := len(matches)
	var ix = 0
	for ix = out.ix; ix+count-1 < len(*out.bytes); ix++ {
		matched := (*out.bytes)[ix : ix+count]
		if slices.Equal(matched, matches) {
			slice := (*out.bytes)[out.ix:ix]
			if !peek {
				out.ix = ix
			}
			return slice, nil
		}
	}
	return nil, errors.New("insufficient bytes in stream")
}

// 🟧 Model inbound 3270 data as a stream
//    "Inbound" data flows from the 3270 ie this code to the application

type InboundDataStream struct {
	bytes []uint8
}

func NewInboundDataStream() *InboundDataStream {
	in := &InboundDataStream{}
	in.bytes = []uint8{}
	return in
}

func (in *InboundDataStream) Put(u8 uint8) []uint8 {
	in.bytes = append(in.bytes, u8)
	return in.bytes
}

func (in *InboundDataStream) Put16(u16 uint16) []uint8 {
	in.bytes = append(in.bytes, uint8(u16>>8))
	in.bytes = append(in.bytes, uint8(u16&0x00ff))
	return in.bytes
}

func (in *InboundDataStream) PutSlice(slice []uint8) []uint8 {
	in.bytes = append(in.bytes, slice...)
	return in.bytes
}

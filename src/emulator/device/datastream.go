package device

import (
	"errors"
	"slices"
)

// 🟧 Model outbound 3270 data as a stream
//    "Outbound" data flows from the application to the 3270 ie this code

type OutboundDataStream struct {
	bytes *[]byte
	ix    int
}

func NewOutboundDataStream(bytes *[]byte) *OutboundDataStream {
	stream := new(OutboundDataStream)
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

func (out *OutboundDataStream) Next() (byte, error) {
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

func (out *OutboundDataStream) NextSlice(count int) ([]byte, error) {
	return out.nextSliceImpl(count, false)
}

func (out *OutboundDataStream) NextSliceUntil(matches []byte) ([]byte, error) {
	return out.nextSliceUntilImpl(matches, false)
}

func (out *OutboundDataStream) Peek() (byte, error) {
	return out.nextImpl(true)
}

func (out *OutboundDataStream) PeekSlice(count int) ([]byte, error) {
	return out.nextSliceImpl(count, true)
}

func (out *OutboundDataStream) PeekSliceUntil(matches []byte) ([]byte, error) {
	return out.nextSliceUntilImpl(matches, true)
}

func (out *OutboundDataStream) Skip(count int) {
	out.ix += count
}

// 👇 Helpers

func (out *OutboundDataStream) nextImpl(peek bool) (byte, error) {
	if out.HasNext() {
		byte := (*out.bytes)[out.ix]
		if !peek {
			out.ix += 1
		}
		return byte, nil
	} else {
		return 0, errors.New("insufficient bytes in stream")
	}
}

func (out *OutboundDataStream) nextSliceImpl(count int, peek bool) ([]byte, error) {
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

func (out *OutboundDataStream) nextSliceUntilImpl(matches []byte, peek bool) ([]byte, error) {
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
	// 👇 if no match, return slice to end
	slice := (*out.bytes)[out.ix:len(*out.bytes)]
	return slice, errors.New("no matches found in stream")
}

// 🟧 Model inbound 3270 data as a stream
//    "Inbound" data flows from the 3270 ie this code to the application

type InboundDataStream struct {
	bytes []byte
}

func NewInboundDataStream() *InboundDataStream {
	in := new(InboundDataStream)
	in.bytes = []byte{}
	return in
}

func (in *InboundDataStream) Put(byte byte) []byte {
	in.bytes = append(in.bytes, byte)
	return in.bytes
}

func (in *InboundDataStream) Put16(u16 uint16) []byte {
	in.bytes = append(in.bytes, byte(u16>>8))
	in.bytes = append(in.bytes, byte(u16&0x00ff))
	return in.bytes
}

func (in *InboundDataStream) PutSlice(slice []byte) []byte {
	in.bytes = append(in.bytes, slice...)
	return in.bytes
}

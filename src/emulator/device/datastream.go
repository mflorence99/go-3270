package device

import (
	"bytes"
	"errors"
)

// 🟧 Model outbound 3270 data as a stream
//    "Outbound" data flows from the application to the 3270 ie this code

type OutboundDataStream struct {
	u8s *[]byte
	ix  int
}

func NewOutboundDataStream(u8s *[]byte) *OutboundDataStream {
	stream := new(OutboundDataStream)
	stream.u8s = u8s
	stream.ix = 0
	return stream
}

func (out *OutboundDataStream) HasEnough(count int) bool {
	return out.ix+count-1 < len(*out.u8s)
}

func (out *OutboundDataStream) HasNext() bool {
	return out.ix < len(*out.u8s)
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
		u8 := (*out.u8s)[out.ix]
		if !peek {
			out.ix += 1
		}
		return u8, nil
	} else {
		return 0, errors.New("insufficient bytes in stream")
	}
}

func (out *OutboundDataStream) nextSliceImpl(count int, peek bool) ([]byte, error) {
	if out.HasEnough(count) {
		end := out.ix + count
		slice := (*out.u8s)[out.ix:end]
		if !peek {
			out.ix = end
		}
		return slice, nil
	} else {
		rem := (*out.u8s)[out.ix:]
		return rem, errors.New("insufficient bytes in stream")
	}
}

func (out *OutboundDataStream) nextSliceUntilImpl(matches []byte, peek bool) ([]byte, error) {
	rem := (*out.u8s)[out.ix:]
	ix := bytes.Index(rem, matches)
	if ix == -1 {
		return rem, errors.New("no matches found in stream")
	} else {
		slice := rem[0:ix]
		if !peek {
			out.ix += ix
		}
		return slice, nil
	}
}

// 🟧 Model inbound 3270 data as a stream
//    "Inbound" data flows from the 3270 ie this code to the application

type InboundDataStream struct {
	u8s []byte
}

func NewInboundDataStream() *InboundDataStream {
	in := new(InboundDataStream)
	in.u8s = []byte{}
	return in
}

func (in *InboundDataStream) Put(u8 byte) []byte {
	in.u8s = append(in.u8s, u8)
	return in.u8s
}

func (in *InboundDataStream) Put16(u16 uint16) []byte {
	in.u8s = append(in.u8s, byte(u16>>8))
	in.u8s = append(in.u8s, byte(u16&0x00ff))
	return in.u8s
}

func (in *InboundDataStream) PutSlice(slice []byte) []byte {
	in.u8s = append(in.u8s, slice...)
	return in.u8s
}

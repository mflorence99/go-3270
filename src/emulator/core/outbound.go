package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 🟧 Outbound (3270 <- app) data stream

type Outbound struct {
	bus   *Bus
	chars []byte
	ix    int
}

// 🟦 Constructor

func NewOutbound(chars []byte, bus *Bus) *Outbound {
	out := new(Outbound)
	out.bus = bus
	out.chars = chars
	out.ix = 0
	return out
}

// 🟦 Public functions

func (out *Outbound) Bytes() []byte {
	return out.chars
}

func (out *Outbound) HasEnough(count int) bool {
	return out.ix+count-1 < len(out.chars)
}

func (out *Outbound) HasNext() bool {
	return out.ix < len(out.chars)
}

func (out *Outbound) Next() (byte, bool) {
	return out.nextImpl(false)
}

func (out *Outbound) Next16() (uint16, bool) {
	return out.next16Impl(false)
}

func (out *Outbound) NextSlice(count int) ([]byte, bool) {
	return out.nextSliceImpl(count, false)
}

func (out *Outbound) NextSliceUntil(matches []byte) ([]byte, bool) {
	return out.nextSliceUntilImpl(matches, false)
}

func (out *Outbound) Peek() (byte, bool) {
	return out.nextImpl(true)
}

func (out *Outbound) Peek16() (uint16, bool) {
	return out.next16Impl(true)
}

func (out *Outbound) PeekSlice(count int) ([]byte, bool) {
	return out.nextSliceImpl(count, true)
}

func (out *Outbound) PeekSliceUntil(matches []byte) ([]byte, bool) {
	return out.nextSliceUntilImpl(matches, true)
}

func (out *Outbound) Rest() []byte {
	rest, _ := out.nextSliceImpl(len(out.chars)-out.ix, false)
	return rest
}

func (out *Outbound) Skip(count int) bool {
	if out.HasEnough(count) {
		out.ix += count
		return true
	}
	return false
}

// 🟦 Helpers

func (out *Outbound) nextImpl(peek bool) (byte, bool) {
	if out.HasNext() {
		byte := out.chars[out.ix]
		if !peek {
			out.ix++
		}
		return byte, true
	} else {
		return 0, false
	}
}

func (out *Outbound) next16Impl(peek bool) (uint16, bool) {
	slice, ok := out.nextSliceImpl(2, peek)
	if !ok {
		return 0, false
	}
	return binary.BigEndian.Uint16(slice), true
}

func (out *Outbound) nextSliceImpl(count int, peek bool) ([]byte, bool) {
	if out.HasEnough(count) {
		end := out.ix + count
		slice := out.chars[out.ix:end]
		if !peek {
			out.ix = end
		}
		return slice, true
	} else {
		rem := out.chars[out.ix:]
		return rem, false
	}
}

func (out *Outbound) nextSliceUntilImpl(matches []byte, peek bool) ([]byte, bool) {
	rem := out.chars[out.ix:]
	ix := bytes.Index(rem, matches)
	if ix == -1 {
		return rem, false
	} else {
		slice := rem[0:ix]
		if !peek {
			out.ix += ix
		}
		return slice, true
	}
}

// 🟦 "Must" functions

func (out *Outbound) MustNext() byte {
	return out.mustNextImpl(false)
}

func (out *Outbound) MustNext16() uint16 {
	u16, ok := out.Next16()
	if !ok {
		out.mustPanic()
	}
	return u16
}

func (out *Outbound) MustNextSlice(count int) []byte {
	return out.mustNextSliceImpl(count, false)
}

func (out *Outbound) MustNextSliceUntil(matches []byte) []byte {
	return out.mustNextSliceUntilImpl(matches, false)
}

func (out *Outbound) MustPeek() byte {
	return out.mustNextImpl(true)
}

func (out *Outbound) MustPeekSlice(count int) []byte {
	return out.mustNextSliceImpl(count, true)
}

func (out *Outbound) MustPeekSliceUntil(matches []byte) []byte {
	return out.mustNextSliceUntilImpl(matches, true)
}

func (out *Outbound) MustSkip(count int) {
	ok := out.Skip(count)
	if !ok {
		out.mustPanic()
	}
}

// 🟦 "Must" Helpers

func (out *Outbound) mustNextImpl(peek bool) byte {
	char, ok := out.nextImpl(peek)
	if !ok {
		out.mustPanic()
	}
	return char
}

func (out *Outbound) mustNextSliceImpl(count int, peek bool) []byte {
	slice, ok := out.nextSliceImpl(count, peek)
	if !ok {
		out.mustPanic()
	}
	return slice
}

func (out *Outbound) mustNextSliceUntilImpl(matches []byte, peek bool) []byte {
	slice, ok := out.nextSliceUntilImpl(matches, peek)
	if !ok {
		out.mustPanic()
	}
	return slice
}

func (out *Outbound) mustPanic() {
	out.bus.PubPanic(fmt.Sprintf("🔥 Internal error: outbound data stream corrupted at offset %d", out.ix))
}

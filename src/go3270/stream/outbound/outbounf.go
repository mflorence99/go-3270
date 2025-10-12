package outbound

import (
	"bytes"
)

// 🔥 "Outbound" data flows from the application to the 3270 ie this code

type Outbound struct {
	bytes *[]byte
	ix    int
}

func New(bytes *[]byte) *Outbound {
	out := new(Outbound)
	out.bytes = bytes
	out.ix = 0
	return out
}

func (out *Outbound) HasEnough(count int) bool {
	return out.ix+count-1 < len(*out.bytes)
}

func (out *Outbound) HasNext() bool {
	return out.ix < len(*out.bytes)
}

func (out *Outbound) Next() (byte, bool) {
	return out.nextImpl(false)
}

func (out *Outbound) Next16() (uint16, bool) {
	hi, ok := out.Next()
	if !ok {
		return 0, false
	}
	lo, ok := out.Next()
	if !ok {
		return 0, false
	}
	return (uint16(hi) * 256) + uint16(lo), true
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

func (out *Outbound) PeekSlice(count int) ([]byte, bool) {
	return out.nextSliceImpl(count, true)
}

func (out *Outbound) PeekSliceUntil(matches []byte) ([]byte, bool) {
	return out.nextSliceUntilImpl(matches, true)
}

func (out *Outbound) Skip(count int) {
	out.ix += count
}

// 👇 Helpers

func (out *Outbound) nextImpl(peek bool) (byte, bool) {
	if out.HasNext() {
		byte := (*out.bytes)[out.ix]
		if !peek {
			out.ix++
		}
		return byte, true
	} else {
		return 0, false
	}
}

func (out *Outbound) nextSliceImpl(count int, peek bool) ([]byte, bool) {
	if out.HasEnough(count) {
		end := out.ix + count
		slice := (*out.bytes)[out.ix:end]
		if !peek {
			out.ix = end
		}
		return slice, true
	} else {
		rem := (*out.bytes)[out.ix:]
		return rem, false
	}
}

func (out *Outbound) nextSliceUntilImpl(matches []byte, peek bool) ([]byte, bool) {
	rem := (*out.bytes)[out.ix:]
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

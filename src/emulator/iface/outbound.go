package iface

type Outbound interface {
	Bytes() []byte
	HasEnough(count int) bool
	HasNext() bool
	MustNext() byte
	MustNext16() uint16
	MustNextSlice(count int) []byte
	MustNextSliceUntil(matches []byte) []byte
	MustPeek() byte
	MustPeekSlice(count int) []byte
	MustPeekSliceUntil(matches []byte) []byte
	MustSkip(count int)
	Next() (byte, bool)
	Next16() (uint16, bool)
	NextSlice(count int) ([]byte, bool)
	NextSliceUntil(matches []byte) ([]byte, bool)
	Peek() (byte, bool)
	PeekSlice(count int) ([]byte, bool)
	PeekSliceUntil(matches []byte) ([]byte, bool)
	Rest() []byte
	Skip(count int) bool
}

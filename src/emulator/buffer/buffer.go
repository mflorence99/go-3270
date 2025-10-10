package buffer

import (
	"emulator/stack"
)

type Buffer struct {
	addr    int
	buffer  []Cell
	changes *stack.Stack[int]
}

func New(size int) *Buffer {
	b := new(Buffer)
	b.buffer = make([]Cell, size)
	return b
}

func (b *Buffer) Seek(addr int) (int, bool) {
	if addr >= len(b.buffer) {
		return 0, false
	}
	b.addr = addr
	return b.addr, true
}

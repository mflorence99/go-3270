package buffer

import (
	"fmt"
	"go3270/emulator/pubsub"
	"go3270/emulator/types"
)

// 🟧 Basic buffer operations

// 🔥 NOTE: the buffer will always hold the original EBCDIC encodings

type Buffer struct {
	addr int
	bus  *pubsub.Bus
	buf  []*Cell
	cfg  types.Config
	mode types.Mode
}

// 🟦 Constructor

func NewBuffer(bus *pubsub.Bus) *Buffer {
	b := new(Buffer)
	b.bus = bus
	// 👇 subscriptions
	b.bus.SubConfig(b.configure)
	b.bus.SubReset(b.reset)
	return b
}

func (b *Buffer) configure(cfg types.Config) {
	b.cfg = cfg
	b.reset()
}

func (b *Buffer) reset() {
	b.buf = make([]*Cell, b.cfg.Cols*b.cfg.Rows)
	b.mode = types.FIELD_MODE
}

// 🟦 Housekeeping functions

//    Addr() get current buffer address
//    Len() get number of cell slots in buffer
//    Mode() reports the buffer's reply mode
//    Peek() cell at given address
//    Replace() cell at given address
//    Seek() reposition buffer address
//    SetMode() sets the buffer's reply mode
//    WrappingSeek() uses a circular buffer address

func (b *Buffer) Addr() int {
	return b.addr
}

func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Mode() types.Mode {
	return b.mode
}

func (b *Buffer) Peek(addr int) (*Cell, bool) {
	if addr < 0 || addr >= len(b.buf) {
		return nil, false
	}
	return b.buf[addr], true
}

func (b *Buffer) Replace(cell *Cell, addr int) (*Cell, bool) {
	if addr < 0 || addr >= len(b.buf) {
		return nil, false
	}
	b.buf[addr] = cell
	return b.buf[addr], true
}

func (b *Buffer) Seek(addr int) (int, bool) {
	if addr < 0 || addr >= len(b.buf) {
		return 0, false
	}
	b.addr = addr
	return b.addr, true
}

func (b *Buffer) SetMode(mode types.Mode) types.Mode {
	if mode > b.mode {
		b.mode = mode
	}
	return b.mode
}

func (b *Buffer) WrappingSeek(addr int) int {
	b.addr = addr % b.Len()
	return b.addr
}

// 🟦 "Must" functions

//    MustPeek() panics if invaliud addr supplied
//    MustReplace() panics if invaliud addr supplied
//    MustSeek() treats addr as a circular ref and wraps

func (b *Buffer) MustPeek(addr int) *Cell {
	cell, ok := b.Peek(addr)
	if !ok {
		b.mustAddr(addr)
	}
	return cell
}

func (b *Buffer) MustReplace(cell *Cell, addr int) *Cell {
	cell, ok := b.Replace(cell, addr)
	if !ok {
		b.mustAddr(addr)
	}
	return cell
}

func (b *Buffer) MustSeek(addr int) int {
	addr, ok := b.Seek(addr)
	if !ok {
		b.mustAddr(addr)
	}
	return addr
}

func (b *Buffer) mustAddr(addr int) {
	row, col := b.cfg.Addr2RC(addr)
	b.bus.PubPanic(fmt.Sprintf("🔥 Internal error: buffer addr %d/%d out of range", row, col))
}

// 🟦 Get functions

//    Get() cell at current address, no side effects
//    GetNext() cell at current address + 1, honoring wrap
//    GetPrev() cell at current address - 1, honoring wrap

func (b *Buffer) Get() (*Cell, int) {
	return b.buf[b.addr], b.addr
}

func (b *Buffer) GetNext() (*Cell, int) {
	addr := b.addr + 1
	if addr >= len(b.buf) {
		addr = 0
	}
	return b.buf[addr], addr
}

func (b *Buffer) PrevGet() (*Cell, int) {
	addr := b.addr - 1
	if addr < 0 {
		addr = len(b.buf) - 1
	}
	return b.buf[addr], addr
}

// 🟦 Set functions

//    Set() cell at current address, no side effects
//    SetAndNext() replace cell at current address, then advance to next
//    PrevAndSet() point to previous cell then replace it

func (b *Buffer) Set(cell *Cell) int {
	b.buf[b.addr] = cell
	return b.addr
}

func (b *Buffer) SetAndNext(cell *Cell) int {
	addr := b.Set(cell)
	if b.addr++; b.addr >= len(b.buf) {
		b.addr = 0
	}
	return addr
}

func (b *Buffer) PrevAndSet(cell *Cell) int {
	if b.addr--; b.addr < 0 {
		b.addr = len(b.buf) - 1
	}
	addr := b.Set(cell)
	return addr
}

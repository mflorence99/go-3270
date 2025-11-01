package buffer

import (
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/pubsub"
)

// 🟧 Basic buffer operations

// 🔥 NOTE: the buffer will always hold ASCII characters

type Buffer struct {
	addr int
	bus  *pubsub.Bus
	buf  []*Cell
	cfg  pubsub.Config
	mode consts.Mode
}

func NewBuffer(bus *pubsub.Bus) *Buffer {
	b := new(Buffer)
	b.bus = bus
	// 👇 subscriptions
	b.bus.SubConfig(b.configure)
	b.bus.SubReset(b.reset)
	return b
}

func (b *Buffer) configure(cfg pubsub.Config) {
	b.cfg = cfg
	b.reset()
}

func (b *Buffer) reset() {
	b.buf = make([]*Cell, b.cfg.Cols*b.cfg.Rows)
	b.mode = consts.FIELD_MODE
}

// 🟦 Housekeeping methods

//    Addr() get current buffer address
//    Len() get number of cell slots in buffer
//    Mode() reports the buffer's reply mode
//    Peek() cell at given address
//    Replace() cell at given address
//    Seek() reposition buffer address
//    SetMode() sets the buffer's reply mode

func (b *Buffer) Addr() int {
	return b.addr
}

func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Mode() consts.Mode {
	return b.mode
}

func (b *Buffer) Peek(addr int) (*Cell, bool) {
	if addr >= len(b.buf) {
		return nil, false
	}
	return b.buf[addr], true
}

func (b *Buffer) Replace(cell *Cell, addr int) {
	b.buf[addr] = cell
}

func (b *Buffer) Seek(addr int) {
	// 🔥 mot wrap around
	b.addr = addr % b.Len()
}

func (b *Buffer) SetMode(mode consts.Mode) consts.Mode {
	if mode > b.mode {
		b.mode = mode
	}
	return b.mode
}

// 🟦 Get methods

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

// 🟦 Set methods

//    Set() cell at current address, no side effects
//    SetAndNext() replace cell at current address, then advance to next
//    StartFld() like SetAndNext(), but for a pre-fab'd SF field
//    PrevAndSet() point to previous cell then replace it

func (b *Buffer) Set(c *Cell) int {
	b.buf[b.addr] = c
	return b.addr
}

func (b *Buffer) SetAndNext(c *Cell) int {
	addr := b.Set(c)
	if b.addr++; b.addr >= len(b.buf) {
		b.addr = 0
	}
	return addr
}

func (b *Buffer) StartFld(a *attrs.Attrs) int {
	c := Cell{
		Attrs:    a,
		Char:     byte(consts.SF),
		FldAddr:  b.addr,
		FldStart: true,
	}
	return b.SetAndNext(&c)
}

func (b *Buffer) PrevAndSet(c *Cell) int {
	if b.addr--; b.addr < 0 {
		b.addr = len(b.buf) - 1
	}
	addr := b.Set(c)
	return addr
}

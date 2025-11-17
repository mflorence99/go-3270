package core

import (
	"emulator/types"
	"fmt"
)

// 🟧 Basic buffer operations

// 🔥 NOTE: the buffer will always hold the original EBCDIC encodings

type Buffer struct {
	addr uint
	buf  []*Cell
	mode types.Mode

	emu *Emulator // 👈 back pouinter to all common components
}

// 🟦 Constructor

func NewBuffer(emu *Emulator) *Buffer {
	b := new(Buffer)
	b.emu = emu
	// 👇 subscriptions
	b.emu.Bus.SubInit(b.init)
	b.emu.Bus.SubReset(b.reset)
	return b
}

func (b *Buffer) init() {
	b.reset()
}

func (b *Buffer) reset() {
	b.buf = make([]*Cell, b.emu.Cfg.Cols*b.emu.Cfg.Rows)
	b.mode = types.FIELD_MODE
}

// 🟦 Low-level functions

//    Addr() get current buffer address
//    Len() get number of cell slots in buffer
//    Mode() reports the buffer's reply mode
//    Peek() cell at given address
//    Replace() cell at given address
//    Seek() reposition buffer address
//    SetMode() sets the buffer's reply mode
//    WrapAddr() computes a circular buffer address
//    WrappingPeek() uses a circular buffer address
//    WrappingSeek() uses a circular buffer address

func (b *Buffer) Addr() uint {
	return b.addr
}

func (b *Buffer) Len() uint {
	return uint(len(b.buf))
}

func (b *Buffer) Mode() types.Mode {
	return b.mode
}

func (b *Buffer) Peek(addr uint) (*Cell, bool) {
	if addr >= b.Len() {
		return nil, false
	}
	return b.buf[addr], true
}

func (b *Buffer) Replace(cell *Cell, addr uint) (*Cell, bool) {
	if addr >= b.Len() {
		return nil, false
	}
	b.buf[addr] = cell
	return b.buf[addr], true
}

func (b *Buffer) Seek(addr uint) (uint, bool) {
	if addr >= b.Len() {
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

func (b *Buffer) WrapAddr(addr uint) uint {
	return addr % b.Len()
}

func (b *Buffer) WrappingPeek(addr uint) (*Cell, uint) {
	temp := b.WrapAddr(addr)
	cell, _ := b.Peek(temp)
	return cell, temp
}

func (b *Buffer) WrappingSeek(addr uint) uint {
	b.addr = b.WrapAddr(addr)
	return b.addr
}

// 🟦 "Must" functions

//    MustPeek() panics if invaliud addr supplied
//    MustReplace() panics if invaliud addr supplied
//    MustSeek() panics if invaliud addr supplied

func (b *Buffer) MustPeek(addr uint) *Cell {
	cell, ok := b.Peek(addr)
	if !ok {
		b.mustAddr(addr)
	}
	return cell
}

func (b *Buffer) MustReplace(cell *Cell, addr uint) *Cell {
	cell, ok := b.Replace(cell, addr)
	if !ok {
		b.mustAddr(addr)
	}
	return cell
}

func (b *Buffer) MustSeek(addr uint) uint {
	addr, ok := b.Seek(addr)
	if !ok {
		b.mustAddr(addr)
	}
	return addr
}

func (b *Buffer) mustAddr(addr uint) {
	row, col := b.emu.Cfg.Addr2RC(addr)
	b.emu.Bus.PubPanic(fmt.Sprintf("🔥 Internal error: buffer addr %d/%d out of range", row, col))
}

// 🟦 Get functions

//    Get() cell at current address, no side effects
//    GetNext() cell at current address + 1, honoring wrap
//    GetPrev() cell at current address - 1, honoring wrap

func (b *Buffer) Get() (*Cell, uint) {
	return b.buf[b.addr], b.addr
}

func (b *Buffer) GetNext() (*Cell, uint) {
	addr := b.addr
	if addr == b.Len()-1 {
		addr = 0
	} else {
		addr++
	}
	return b.buf[addr], addr
}

func (b *Buffer) PrevGet() (*Cell, uint) {
	addr := b.addr
	if addr == 0 {
		addr = b.Len() - 1
	} else {
		addr--
	}
	return b.buf[addr], addr
}

// 🟦 Set functions

//    Set() cell at current address, no side effects
//    SetAndNext() replace cell at current address, then advance to next
//    PrevAndSet() point to previous cell then replace it

func (b *Buffer) Set(cell *Cell) uint {
	b.buf[b.addr] = cell
	return b.addr
}

func (b *Buffer) SetAndNext(cell *Cell) uint {
	addr := b.Set(cell)
	if b.addr == b.Len()-1 {
		b.addr = 0
	} else {
		b.addr++
	}
	return addr
}

func (b *Buffer) PrevAndSet(cell *Cell) uint {
	if b.addr == 0 {
		b.addr = b.Len() - 1
	} else {
		b.addr--
	}
	addr := b.Set(cell)
	return addr
}

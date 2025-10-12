package buffer

import (
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/stack"
	"strings"
)

// 🔥 NOTE: the buffer will always hold ASCII characters

type Buffer struct {
	Changes *stack.Stack[int]

	addr   int
	buffer []*Cell
}

func NewBuffer(size int) *Buffer {
	b := new(Buffer)
	b.buffer = make([]*Cell, size)
	b.Changes = stack.NewStack[int](1)
	b.Erase()
	return b
}

// 🟦 Housekeeping methods

//    Erase() fill the buffer with protected cells
//    Len() get number of cell slots in buffer
//    Peek() cell at given address
//    Seek() reposition buffer address

func (b *Buffer) Erase() {
	for ix := range b.buffer {
		b.buffer[ix] = &Cell{Attrs: &attrs.Attrs{Protected: true}}
	}
	for !b.Changes.Empty() {
		b.Changes.Pop()
	}
}

func (b *Buffer) Len() int {
	return len(b.buffer)
}

func (b *Buffer) Peek(addr int) (*Cell, bool) {
	if addr >= len(b.buffer) {
		return nil, false
	}
	return b.buffer[addr], true
}

func (b *Buffer) Seek(addr int) (int, bool) {
	if addr >= len(b.buffer) {
		return -1, false
	}
	b.addr = addr
	return b.addr, true
}

// 🟦 Get methods

//    Get() cell at current address, no side effects
//    GetNext() cell at current address + 1, honoring wrap
//    GetPrev() cell at current address - 1, honoring wrap

func (b *Buffer) Get() (*Cell, int) {
	return b.buffer[b.addr], b.addr
}

func (b *Buffer) GetNext() (*Cell, int) {
	addr := b.addr + 1
	if addr >= len(b.buffer) {
		addr = 0
	}
	return b.buffer[addr], addr
}

func (b *Buffer) PrevGet() (*Cell, int) {
	addr := b.addr - 1
	if addr < 0 {
		addr = len(b.buffer) - 1
	}
	return b.buffer[addr], addr
}

// 🟦 Set methods

//    Set() cell at current address, no side effects
//    SetAndNext() replace cell at current address, then advance to next
//    StartFld() like SetAndNext(), but for a pre-fab'd SF field
//    PrevAndSet() point to previous cell then replace it

func (b *Buffer) Set(c *Cell) int {
	b.buffer[b.addr] = c
	b.Changes.Push(b.addr)
	return b.addr
}

func (b *Buffer) SetAndNext(c *Cell) int {
	addr := b.Set(c)
	if b.addr++; b.addr >= len(b.buffer) {
		b.addr = 0
	}
	return addr
}

func (b *Buffer) StartFld(attrs *attrs.Attrs) int {
	c := Cell{
		Attrs:    attrs,
		Char:     byte(consts.SF),
		FldStart: true,
	}
	return b.SetAndNext(&c)
}

func (b *Buffer) PrevAndSet(c *Cell) int {
	if b.addr--; b.addr < 0 {
		b.addr = len(b.buffer) - 1
	}
	addr := b.Set(c)
	return addr
}

// 🟦 Keystroke methods

//    Keyin() update char in current cell, then advance to next
//    Backspace() point to previous cell then update its char

func (b *Buffer) Keyin(char byte) (int, bool) {
	c, _ := b.Get()
	// 👇 validate data entry into current cell
	numlock := c.Attrs.Numeric && !strings.Contains("0123456789.", string(char))
	prot := c.Attrs.Hidden || c.Attrs.Protected
	if numlock || prot {
		return -1, false
	}
	// 👇 update cell and advance to next
	c.Char = char
	c.Attrs.Modified = true
	addr := b.SetAndNext(c)
	// 👇 set the MDT flag at the field level
	ok := b.setFldMDT(c.FldAddr)
	if !ok {
		return -1, false
	}
	return addr, true
}

func (b *Buffer) Backspace() (int, bool) {
	// 👇-validate data entry into previous cell
	c, addr := b.PrevGet()
	prot := c.Attrs.Hidden || c.Attrs.Protected
	if prot {
		return -1, false
	}
	// 👇 reposition to previous cell and update it
	b.Seek(addr)
	c.Char = 0x00
	c.Attrs.Modified = true
	addr = b.Set(c)
	// 👇 set the MDT flag at the field level
	ok := b.setFldMDT(c.FldAddr)
	if !ok {
		return -1, false
	}
	return addr, true
}

func (b *Buffer) setFldMDT(fldAddr int) bool {
	fld, ok := b.Peek(fldAddr)
	if !ok {
		return false
	}
	fld.Attrs.Modified = true
	return true
}

package buffer

import (
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"strings"
)

// 🔥 NOTE: the buffer will always hold ASCII characters

type Buffer struct {
	addr   int
	bus    *pubsub.Bus
	buf    []*Cell
	cfg    pubsub.Config
	deltas *utils.Stack[int]
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
	b.buf = make([]*Cell, cfg.Cols*cfg.Rows)
	b.deltas = utils.NewStack[int](1)
}

func (b *Buffer) reset() {
	b.Erase(0x00)
}

// 🟦 Housekeeping methods

//    Addr() get current buffer address
//    Chars() extracts ASCII chars from buffer for debugging
//    Deltas() returns stack of changes
//    Flds() returns slice of all fields
//    Len() get number of cell slots in buffer
//    Peek() cell at given address
//    Seek() reposition buffer address

func (b *Buffer) Addr() int {
	return b.addr
}

func (b *Buffer) Chars() []byte {
	chars := make([]byte, b.Len())
	for ix, cell := range b.buf {
		if cell != nil {
			chars[ix] = cell.Char
		}
	}
	return chars
}

func (b *Buffer) Deltas() *utils.Stack[int] {
	return b.deltas
}

func (b *Buffer) Erase(char byte) {
	b.addr = 0
	for ix := range b.buf {
		b.buf[ix] = &Cell{Attrs: &attrs.Attrs{Protected: true}, Char: char}
	}
	for !b.deltas.Empty() {
		b.deltas.Pop()
	}
}

func (b *Buffer) Flds() [][]*Cell {
	fld := make([]*Cell, 0)
	flds := make([][]*Cell, 0)
	for _, cell := range b.buf {
		if cell != nil && cell.FldStart {
			if len(fld) > 0 {
				flds = append(flds, fld)
				fld = make([]*Cell, 0)
			}
		}
		fld = append(fld, cell)
	}
	// 👇 don't forget the last field
	if len(fld) > 0 {
		flds = append(flds, fld)
	}
	return flds
}

func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Peek(addr int) (*Cell, bool) {
	if addr >= len(b.buf) {
		return nil, false
	}
	return b.buf[addr], true
}

func (b *Buffer) Seek(addr int) (int, bool) {
	if addr >= len(b.buf) {
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
	b.deltas.Push(b.addr)
	return b.addr
}

func (b *Buffer) SetAndNext(c *Cell) int {
	addr := b.Set(c)
	if b.addr++; b.addr >= len(b.buf) {
		b.addr = 0
	}
	return addr
}

func (b *Buffer) StartFld(attrs *attrs.Attrs) int {
	c := Cell{
		Attrs:    attrs,
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

// 🟦 Keystroke methods

//    Keyin() updates char in current cell, then advances to next
//    Backspace() points to previous cell then updates its char
//    Tab() skips forward or backward to the next unprotected field

func (b *Buffer) Keyin(char byte) (int, bool) {
	c, _ := b.Get()
	// 👇 validate data entry into current cell
	numlock := c.Attrs.Numeric && !strings.Contains("-0123456789.", string(char))
	prot := c.FldStart || c.Attrs.Hidden || c.Attrs.Protected
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
	// 👇 validate data entry into previous cell
	c, addr := b.PrevGet()
	prot := c.FldStart || c.Attrs.Hidden || c.Attrs.Protected
	if prot {
		return -1, false
	}
	// 👇 reposition to previous cell and update it
	b.Seek(addr)
	c.Char = ' '
	c.Attrs.Modified = true
	addr = b.Set(c)
	// 👇 set the MDT flag at the field level
	ok := b.setFldMDT(c.FldAddr)
	if !ok {
		return -1, false
	}
	return addr, true
}

func (b *Buffer) Tab(dir int) (int, bool) {
	start := b.addr
	addr := b.addr
	for ix := 0; ; ix++ {
		// 👇 wrap to the start means no unprotected field
		if addr == start && ix > 0 {
			break
		}
		// 👇 look at the "next" cell
		addr += dir
		if addr < 0 {
			addr = len(b.buf) - 1
		} else if addr >= len(b.buf) {
			addr = 0
		}
		// 👇 see if we've hit an unprotected field start
		cell, _ := b.Peek(addr)
		if cell.FldStart && !cell.Attrs.Protected {
			b.Seek(addr) // 👈 go to FldStart
			_, addr := b.GetNext()
			b.Seek(addr) // 👈 now to first char
			return addr, true
		}
	}
	return -1, false
}

func (b *Buffer) setFldMDT(fldAddr int) bool {
	fld, ok := b.Peek(fldAddr)
	if !ok {
		return false
	}
	fld.Attrs.Modified = true
	return true
}

package buffer

import (
	"fmt"
	"go3270/emulator/consts"
	"go3270/emulator/pubsub"
)

// 🟧 View the buffer as an array of cells

type Cells struct {
	buf *Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
}

// 🟦 Constructor

func NewCells(bus *pubsub.Bus, buf *Buffer) *Cells {
	c := new(Cells)
	c.buf = buf
	c.bus = bus
	// 👇 subscriptions
	c.bus.SubConfig(c.configure)
	c.bus.SubReset(c.reset)
	return c
}

func (c *Cells) configure(cfg pubsub.Config) {
	c.cfg = cfg
	c.reset()
}

func (c *Cells) reset() {
	for addr := 0; addr < c.buf.Len(); addr++ {
		cell, _ := c.buf.Peek(addr)
		if cell == nil {
			cell = NewCell()
			c.buf.Replace(cell, addr)
		}
	}
}

// 🟦 Public functions

func (c *Cells) FillOlder2Left(cell *Cell, stop int) {
	for ix := stop; ; ix-- {
		// 🔥 note wrap around
		addr := (ix - 1) % c.buf.Len()
		left, ok := c.buf.Peek(addr)
		if ok && left.FldGen != cell.FldGen {
			left.Attrs = cell.Attrs
			left.Char = 0x00
			left.FldAddr = cell.FldAddr
			left.FldStart = false
			left.FldEnd = false
			left.FldGen = cell.FldGen
		} else {
			break
		}
	}
}

// 🟦 Public order-based functions

func (c *Cells) EUA(start, stop int) bool {
	if stop < c.buf.Len() {
		for addr := c.buf.Seek(start); addr != stop; addr = c.buf.Seek(addr + 1) {
			cell, _ := c.buf.Get()
			if !cell.Attrs.Protected {
				// TODO 🔥 spec says to reset any character attributes ie revert to fld but this just makes everything blank
				sf, _ := c.buf.Peek(cell.FldAddr)
				cell.Attrs = sf.Attrs
				cell.Char = 0x00
				c.buf.Replace(cell, addr)
			}
		}
		return true
	} else {
		println(fmt.Sprintf("🔥 Inavlid stop address %d in EUA order terminates write", stop))
		return false
	}
}

func (c *Cells) MF(chars []byte) {
	cell, _ := c.buf.Get()
	cell.Attrs = consts.NewModifiedAttrs(cell.Attrs, chars)
	c.buf.SetAndNext(cell)
}

func (c *Cells) RA(cell *Cell, start, stop int) bool {
	if stop < c.buf.Len() {
		for addr := c.buf.Seek(start); addr != stop; addr = c.buf.Seek(addr + 1) {
			copy := *cell
			c.buf.Replace(&copy, addr)
		}
		return true
	} else {
		println(fmt.Sprintf("🔥 Inavlid stop address %d in RA order terminates write", stop))
		return false
	}
}

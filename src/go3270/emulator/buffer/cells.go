package buffer

import (
	"fmt"
	"go3270/emulator/pubsub"
	"go3270/emulator/types"
)

// 🟧 View the buffer as an array of cells

type Cells struct {
	buf *Buffer
	bus *pubsub.Bus
	cfg types.Config
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

func (c *Cells) configure(cfg types.Config) {
	c.cfg = cfg
	c.reset()
}

func (c *Cells) reset() {
	for addr := 0; addr < c.buf.Len(); addr++ {
		cell := c.buf.MustPeek(addr)
		if cell == nil {
			cell = NewCell()
			c.buf.MustReplace(cell, addr)
		}
	}
}

// 🟦 Public functions

func (c *Cells) EUA(start, stop int) bool {
	if stop < c.buf.Len() {
		addr := c.buf.WrappingSeek(start)
		for addr != stop {
			cell, _ := c.buf.Get()
			if !cell.Attrs.Protected {
				sf := c.buf.MustPeek(cell.FldAddr)
				cell.Attrs = sf.Attrs
				cell.Char = 0x00
				c.buf.MustReplace(cell, addr)
			}
			addr = c.buf.WrappingSeek(addr + 1)
		}
		return true
	} else {
		println(fmt.Sprintf("🔥 Invalid stop address %d in EUA order terminates write", stop))
		return false
	}
}

func (c *Cells) MF(chars []byte) {
	cell, _ := c.buf.Get()
	cell.Attrs = types.NewModifiedAttrs(cell.Attrs, chars)
	c.buf.SetAndNext(cell)
}

func (c *Cells) RA(cell *Cell, start, stop int) bool {
	if stop < c.buf.Len() {
		addr := c.buf.WrappingSeek(start)
		for addr != stop {
			copy := *cell
			c.buf.MustReplace(&copy, addr)
			addr = c.buf.WrappingSeek(addr + 1)
		}
		return true
	} else {
		println(fmt.Sprintf("🔥 Invalid stop address %d in RA order terminates write", stop))
		return false
	}
}

func (c *Cells) RB() []byte {
	chars := make([]byte, 0)
	mode := c.buf.Mode()
	var fldAttrs []byte
	for addr := 0; addr < c.buf.Len(); addr++ {
		cell := c.buf.MustPeek(addr)
		if cell.FldStart {
			if mode == types.FIELD_MODE {
				chars = append(chars, byte(types.SF))
				chars = append(chars, cell.Attrs.Byte())
			} else {
				chars = append(chars, byte(types.SFE))
				fldAttrs = cell.Attrs.Bytes()
				chars = append(chars, byte(len(fldAttrs)/2))
				chars = append(chars, fldAttrs...)
			}
		} else if cell.Attrs.CharAttr {
			charAttrs := cell.Attrs.Bytes()
			for ix := 0; ix < len(charAttrs); ix += 2 {
				if len(charAttrs) != len(fldAttrs) || charAttrs[ix] != fldAttrs[ix] {
					chars = append(chars, byte(types.SA))
					chars = append(chars, charAttrs[ix])
					chars = append(chars, charAttrs[ix+1])
				}
			}
			fldAttrs = charAttrs
			chars = append(chars, cell.Char)
		} else {
			chars = append(chars, cell.Char)
		}
	}
	return chars
}

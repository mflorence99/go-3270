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
	var fldAttrs *types.Attrs
	for addr := 0; addr < c.buf.Len(); addr++ {
		cell := c.buf.MustPeek(addr)
		switch {

		// 👇 delineate FldStart with SF/SFE orders
		case cell.FldStart:
			if mode == types.FIELD_MODE {
				chars = append(chars, byte(types.SF))
				chars = append(chars, cell.Attrs.Byte())
			} else {
				chars = append(chars, byte(types.SFE))
				fldAttrs = cell.Attrs
				raw := fldAttrs.Bytes()
				chars = append(chars, byte(len(raw)/2))
				chars = append(chars, raw...)
			}

		// 👇 emit SA everytime attribute changes
		case cell.Attrs.CharAttr:
			charAttrs := cell.Attrs
			delta := charAttrs.Diff(fldAttrs)
			raw := delta.Bytes()
			for ix := 0; ix < len(raw); ix += 2 {
				chars = append(chars, byte(types.SA))
				chars = append(chars, raw[ix])
				chars = append(chars, raw[ix+1])
			}
			fldAttrs = charAttrs
			chars = append(chars, cell.Char)

		// 👇 just the data
		default:
			chars = append(chars, cell.Char)

		}
	}
	return chars
}

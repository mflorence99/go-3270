package core

import (
	"emulator/types"
)

// 🟧 View the buffer as an array of cells

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

type Cells struct {
	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewCells(emu *Emulator) *Cells {
	c := new(Cells)
	c.emu = emu
	// 👇 subscriptions
	c.emu.Bus.SubInit(c.init)
	c.emu.Bus.SubReset(c.reset)
	return c
}

func (c *Cells) init() {
	c.reset()
}

func (c *Cells) reset() {
	for addr := uint(0); addr < c.emu.Buf.Len(); addr++ {
		cell := c.emu.Buf.MustPeek(addr)
		if cell == nil {
			cell = NewCell(c.emu)
			c.emu.Buf.MustReplace(cell, addr)
		}
	}
}

// 🟦 Public functions

// 👁️ Erase Unprotected to Address order pp 4-10 to 4-11
func (c *Cells) EUA(start, stop uint) {
	addr := c.emu.Buf.WrappingSeek(start)
	for addr != stop {
		cell, _ := c.emu.Buf.Get()
		if !cell.Attrs.Protected {
			sf, ok := cell.GetFldStart()
			if ok {
				cell.Attrs = sf.Attrs
			}
			cell.Char = 0x00
			c.emu.Buf.MustReplace(cell, addr)
		}
		addr = c.emu.Buf.WrappingSeek(addr + 1)
	}
}

// 👁️ Repeat to Address order pp 4-9 to 4-10
func (c *Cells) RA(cell *Cell, start, stop uint) {
	addr := c.emu.Buf.WrappingSeek(start)
	for addr != stop {
		copy := *cell
		c.emu.Buf.MustReplace(&copy, addr)
		addr = c.emu.Buf.WrappingSeek(addr + 1)
	}
}

// 👁️ Read Buffer command pp 3-12 to 3-13
func (c *Cells) RB() []byte {
	chars := make([]byte, 0)
	mode := c.emu.Buf.Mode()
	var fldAttrs *types.Attrs
	for addr := uint(0); addr < c.emu.Buf.Len(); addr++ {
		cell := c.emu.Buf.MustPeek(addr)
		switch {

		// 👇 delineate FldStart with SF/SFE orders
		case cell.IsFldStart():
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
			chars = append(chars, cell.Char)
			// 👇 now the char attrs take over
			fldAttrs = charAttrs

		// 👇 just the data
		default:
			chars = append(chars, cell.Char)

		}
	}
	return chars
}

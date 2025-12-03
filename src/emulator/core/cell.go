package core

import "emulator/types"

// 🟧 Cell in buffer

type Cell struct {
	Attrs *types.Attrs
	Char  byte

	emu     *Emulator // 👈 back pointer to all common components
	fldAddr uint      // 👈 buffer address of Fld
	inFld   bool      // 👈 true if this cell is in a field
}

// 🟦 Internal constructor 👁️ cells.go

func NewCell(emu *Emulator) *Cell {
	c := new(Cell)
	c.Attrs = types.NewDefaultAttrs()
	c.emu = emu
	return c
}

// 🟦 Public functions

func (c *Cell) FindFld() (*Fld, bool) {
	fldAddr, ok := c.GetFldAddr()
	if !ok {
		return nil, false
	}
	return c.emu.Flds.FindFld(fldAddr)
}

func (c *Cell) GetFldAddr() (uint, bool) {
	if !c.inFld {
		return 0, false
	}
	return c.fldAddr, true
}

func (c *Cell) GetFldHome() (*Cell, bool) {
	addr, ok := c.GetFldAddr()
	if !ok {
		return nil, false
	}
	home, _ := c.emu.Buf.WrappingPeek(int(addr) + 1)
	// 🔥 a Fld with no Cells has no home Cell
	if home.IsFldStart() {
		return nil, false
	}
	return home, true
}

func (c *Cell) GetFldStart() (*Cell, bool) {
	addr, ok := c.GetFldAddr()
	if !ok {
		return nil, false
	}
	sf := c.emu.Buf.MustPeek(addr)
	return sf, true
}

func (c *Cell) IsFldHome() bool {
	home, ok := c.GetFldHome()
	if !ok {
		return false
	}
	return c == home
}

func (c *Cell) IsFldStart() bool {
	order := types.Order(c.Char)
	return order == types.SF || order == types.SFE
}

func (c *Cell) SetFldAddr(fldAddr uint) {
	c.fldAddr = fldAddr
	c.inFld = true
}

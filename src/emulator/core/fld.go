package core

import (
	"emulator/conv"
	"strings"
)

// 🟧 Field of cells in buffer

type Fld struct {
	Cells []*Cell

	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Internal constructor 👁️ flds.go

func NewFld(sf *Cell, emu *Emulator) *Fld {
	f := new(Fld)
	f.Cells = make([]*Cell, 1)
	f.Cells[0] = sf
	f.emu = emu
	return f
}

// 🟦 Stringer implementation

func (f Fld) String() string {
	var b strings.Builder
	for ix := 1; ix < len(f.Cells); ix++ {
		cell := f.Cells[ix]
		if cell.Char >= 0x40 {
			b.WriteRune(conv.E2Rune(cell.Attrs.LCID, cell.Char))
		}
	}
	return strings.TrimSpace(b.String())
}

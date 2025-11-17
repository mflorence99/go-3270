package core

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

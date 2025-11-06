package buffer

import (
	"fmt"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
)

// 🟧 View the buffer as an array of fields

type Flds struct {
	buf  *Buffer
	bus  *pubsub.Bus
	cfg  pubsub.Config
	flds []Fld
}

// 🟦 Constructor

func NewFlds(bus *pubsub.Bus, buf *Buffer) *Flds {
	f := new(Flds)
	f.buf = buf
	f.bus = bus
	// 👇 subscriptions
	f.bus.SubConfig(f.configure)
	f.bus.SubReset(f.reset)
	return f
}

func (f *Flds) configure(cfg pubsub.Config) {
	f.cfg = cfg
	f.reset()
}

func (f *Flds) reset() {
	f.flds = make([]Fld, 0)
}

// 🟦 Builder functions

func (f *Flds) Build(fldGen int) {
	flds := f.buildInitialIndex()
	flds = f.elideOldFlds(flds, fldGen)
	f.flds = f.allThoseCellsAreMine(flds, fldGen)
}

func (f *Flds) buildInitialIndex() []Fld {
	flds := make([]Fld, 0)
	for ix := 0; ix < f.buf.Len(); ix++ {
		cell, _ := f.buf.Peek(ix)
		if cell.FldStart {
			fld := make(Fld, 1)
			fld[0] = cell
			flds = append(flds, fld)
		}
	}
	return flds
}

func (f *Flds) elideOldFlds(flds []Fld, fldGen int) []Fld {
	temp := make([]Fld, 0)
	for ix, fld := range flds {
		switch {
		// 👇 we'll keep completely new fields
		case fld[0].FldGen == fldGen:
			temp = append(temp, fld)
		// 👇 we'll keep fields at the edges
		case ix == 0 || ix == (len(flds)-1):
			temp = append(temp, fld)
		// 👇 we'll keep old fields not sandwiched by new ones
		default:
			var newBefore bool
			for iy := ix - 1; iy >= 0; iy-- {
				if flds[iy][0].FldGen == fldGen {
					newBefore = true
					break
				}
			}
			var newAfter bool
			for iy := ix + 1; iy < len(flds); iy++ {
				if flds[iy][0].FldGen == fldGen {
					newAfter = true
					break
				}
			}
			// 👇 here we keep the field, or lose it
			if !newBefore || !newAfter {
				temp = append(temp, fld)
			} else {
				fld[0].FldStart = false
			}
		}
	}
	return temp
}

func (f *Flds) allThoseCellsAreMine(flds []Fld, fldGen int) []Fld {
	temp := make([]Fld, 0)
	for ix, fld := range flds {
		next := flds[(ix+1)%len(flds)]
		stop := next[0].FldAddr
		// 👇 prepare field start
		sf := fld[0]
		start := sf.FldAddr
		sf.FldGen = fldGen
		for iy := start; ; iy++ {
			// 🔥 note wrap around
			addr := (iy + 1) % f.buf.Len()
			if addr == stop {
				break
			}
			cell, _ := f.buf.Peek(addr)
			// 👇 these are cells that got filled outside of a known  field
			if cell.FldAddr == -1 || cell.Attrs.Default {
				cell.Attrs = sf.Attrs
				cell.FldAddr = start
			}
			// 👇 these are cells from a field that no longer exists
			old, _ := f.buf.Peek(cell.FldAddr)
			if !old.FldStart {
				row, col := f.cfg.Addr2RC(addr)
				println(fmt.Sprintf("Erasing cell %d/%d", row, col))
				cell.Attrs = sf.Attrs
				cell.Char = 0x00
			}
			// 👇 make the cell mine
			cell.FldAddr = start
			cell.FldStart = false
			cell.FldEnd = false
			cell.FldGen = fldGen
			fld = append(fld, cell)
		}
		// 👇 mark field end
		ef := fld[len(fld)-1]
		ef.FldEnd = true
		temp = append(temp, fld)
	}
	return temp
}

// 🟦 Housekeeping functions

func (f *Flds) Get() []Fld {
	return f.flds
}

// 🟦 Public command-based functions

func (f *Flds) EAU() int {
	addr := -1
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		if ok {
			sf.Attrs.Modified = false
			if !sf.Attrs.Protected {
				// 👇 capture address of first unprotected field
				if addr == -1 {
					addr = sf.FldAddr
				}
				// 🔥 reset char and any character attributes
				for ix := 1; ix < len(fld); ix++ {
					cell := fld[ix]
					cell.Char = 0x00
					cell.Attrs = sf.Attrs
				}
			}
		}
	}
	return addr
}

// TODO 🔥 *only* FIELD_MODE *not* coded
func (f *Flds) ReadBuffer() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.flds {
		sf, _ := fld.FldStart()
		chars = append(chars, byte(consts.SF))
		chars = append(chars, sf.Attrs.Byte())
		for ix := 1; ix < len(fld); ix++ {
			cell := fld[ix]
			char := cell.Char
			if char != 0x00 {
				chars = append(chars, conv.A2E(char))
			}
		}
	}
	return chars
}

// TODO 🔥 CHARACTER_MODE *not* coded
func (f *Flds) ReadMDTs() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.flds {
		sf, _ := fld.FldStart()
		if sf.Attrs.Modified {
			chars = append(chars, byte(consts.SBA))
			chars = append(chars, conv.Addr2Bytes(sf.FldAddr+1)...)
			for ix := 1; ix < len(fld); ix++ {
				cell := fld[ix]
				char := cell.Char
				if char != 0x00 {
					chars = append(chars, conv.A2E(char))
				}
			}
		}
	}
	return chars
}

func (f *Flds) ResetMDTs() {
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		if ok {
			sf.Attrs.Modified = false
		}
	}
}

func (f *Flds) SetMDT(cell *Cell) bool {
	fld, ok := f.buf.Peek(cell.FldAddr)
	if !fld.FldStart || !ok {
		return false
	}
	fld.Attrs.Modified = true
	return true
}

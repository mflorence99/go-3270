package core

import (
	"emulator/conv"
	"emulator/types"

	"cmp"
	"slices"
)

// 🟧 View the buffer as an array of fields

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

type Flds struct {
	Flds []*Fld

	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewFlds(emu *Emulator) *Flds {
	f := new(Flds)
	f.emu = emu
	// 👇 subscriptions
	f.emu.Bus.SubInit(f.init)
	f.emu.Bus.SubRender(f.build)
	f.emu.Bus.SubReset(f.reset)
	return f
}

func (f *Flds) init() {
	f.reset()
}

func (f *Flds) reset() {
	f.Flds = make([]*Fld, 0)
}

// 🟦 Builder functions

func (f *Flds) build() {
	flds := f.buildIndex()
	f.Flds = f.buildCells(flds)
}

// 👇 just build the start point of each field
func (f *Flds) buildIndex() []*Fld {
	flds := make([]*Fld, 0)
	for addr := uint(0); addr < f.emu.Buf.Len(); addr++ {
		sf := f.emu.Buf.MustPeek(addr)
		if sf.IsFldStart() {
			fld := NewFld(sf, f.emu)
			flds = append(flds, fld)
		}
	}
	return flds
}

// 👇 look at pairs of Flds, this and the next,
// and group all the Cells in between to this Fld
func (f *Flds) buildCells(flds []*Fld) []*Fld {
	temp := make([]*Fld, 0)
	for ix, fld := range flds {
		// 🔥 note wrap around
		next := flds[(ix+1)%len(flds)]
		sf := fld.Cells[0]
		start, _ := sf.GetFldAddr()
		stop, _ := next.Cells[0].GetFldAddr()
		cell, addr := f.emu.Buf.WrappingPeek(int(start) + 1)
		for addr != stop {
			// 👇 use the field attributes for cells that were never
			//    initialized, or which have (potentially) another field's
			//    attributes, or those that used to belong to a
			//    now-overwritten field
			oldFldAddr, ok := cell.GetFldAddr()
			differentFld := start != oldFldAddr || !ok
			if cell.Attrs.Default || !cell.Attrs.CharAttr || differentFld {
				cell.Attrs = sf.Attrs
			}
			// 👇 make the cell mine
			cell.SetFldAddr(start)
			fld.Cells = append(fld.Cells, cell)
			// 👇 on to the next
			cell, addr = f.emu.Buf.WrappingPeek(int(addr) + 1)
		}
		temp = append(temp, fld)
	}
	return temp
}

// 🟦 Public functions

// 👁️ Erase All Unprotected command p 3-8
func (f *Flds) EAU() (uint, bool) {
	var addr uint
	var firstFld bool
	for _, fld := range f.Flds {
		sf := fld.Cells[0]
		if !sf.Attrs.Protected {
			sf.Attrs.MDT = false
			// 👇 capture address of first unprotected field
			if !firstFld {
				addr, firstFld = sf.GetFldAddr()
			}
			// 🔥 reset char and any character attributes
			for ix := 1; ix < len(fld.Cells); ix++ {
				cell := fld.Cells[ix]
				cell.Char = 0x00
				cell.Attrs = sf.Attrs
			}
		}
	}
	return addr, firstFld
}

func (f *Flds) FindFld(fldAddr uint) (*Fld, bool) {
	ix, ok := slices.BinarySearchFunc(f.Flds, fldAddr,
		func(fld *Fld, fldAddr uint) int {
			addr, _ := fld.Cells[0].GetFldAddr()
			return cmp.Compare(addr, fldAddr)
		})
	if !ok {
		return nil, false
	}
	return f.Flds[ix], true
}

// 👁️ Read Modified command pp 3-13 to 3-15
func (f *Flds) RM() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.Flds {
		sf := fld.Cells[0]
		// 👇 for each changed field
		if sf.Attrs.MDT {
			chars = append(chars, byte(types.SBA))
			addr, _ := sf.GetFldAddr()
			next := f.emu.Buf.WrapAddr(int(addr) + 1)
			chars = append(chars, conv.Addr2Bytes(next)...)
			// 👇 now for each cell in that field
			for ix := 1; ix < len(fld.Cells); ix++ {
				cell := fld.Cells[ix]
				// TODO 🔥 this seems to blow the 1 RFE page input in TSOAPPLS
				// 👇 emit SA order for char attrs different to fld attrs
				// if cell.Attrs.CharAttr {
				// 	delta := cell.Attrs.Diff(sf.Attrs)
				// 	raw := delta.Bytes()
				// 	for ix := 0; ix < len(raw); ix += 2 {
				// 		chars = append(chars, byte(types.SA))
				// 		chars = append(chars, raw[ix])
				// 		chars = append(chars, raw[ix+1])
				// 	}
				// }
				// 👇 suppress null characters
				if cell.Char != 0x00 {
					chars = append(chars, cell.Char)
				}
			}
		}
	}
	return chars
}

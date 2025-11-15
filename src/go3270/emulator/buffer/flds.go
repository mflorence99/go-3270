package buffer

import (
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/types"
)

// 🟧 View the buffer as an array of fields

type Flds struct {
	buf  *Buffer
	bus  *pubsub.Bus
	cfg  types.Config
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

func (f *Flds) configure(cfg types.Config) {
	f.cfg = cfg
	f.reset()
}

func (f *Flds) reset() {
	f.flds = make([]Fld, 0)
}

// 🟦 Builder functions

func (f *Flds) Build() {
	flds := f.buildInitialIndex()
	f.flds = f.allThoseCellsAreMine(flds)
}

// 👇 just the start point of each field
func (f *Flds) buildInitialIndex() []Fld {
	flds := make([]Fld, 0)
	for ix := 0; ix < f.buf.Len(); ix++ {
		cell := f.buf.MustPeek(ix)
		if cell.FldStart {
			fld := make(Fld, 1)
			fld[0] = cell
			flds = append(flds, fld)
		}
	}
	return flds
}

// 👇 look at pairs of fields, this one and the next and assign
//
//	all the cells in between to this field
func (f *Flds) allThoseCellsAreMine(flds []Fld) []Fld {
	temp := make([]Fld, 0)
	for ix, fld := range flds {
		next := flds[(ix+1)%len(flds)]
		stop := next[0].FldAddr
		// 👇 prepare field start
		sf := fld[0]
		start := sf.FldAddr
		for iy := start; ; iy++ {
			// 🔥 note wrap around
			addr := (iy + 1) % f.buf.Len()
			if addr == stop {
				break
			}
			cell := f.buf.MustPeek(addr)
			// 👇 use the field attributes for cells that were never
			//    initialized, or which have (potentially) another field's
			//    attributes, or those that used to belong to a
			//    now-overwritten field,
			if cell.Attrs.Default || !cell.Attrs.CharAttr || sf.FldAddr != cell.FldAddr {
				cell.Attrs = sf.Attrs
			}
			// 👇 make the cell mine
			cell.FldAddr = start
			cell.FldStart = false
			cell.FldEnd = false
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

// 🟦 Public functions

func (f *Flds) EAU() int {
	addr := -1
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		if ok {
			sf.Attrs.MDT = false
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

func (f *Flds) RM() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		// 👇 for each changed field
		if ok && sf.Attrs.MDT {
			chars = append(chars, byte(types.SBA))
			chars = append(chars, conv.Addr2Bytes(sf.FldAddr+1)...)
			// 👇 now for each cell in that field
			for ix := 1; ix < len(fld); ix++ {
				cell := fld[ix]
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
				char := cell.Char
				if char != 0x00 {
					chars = append(chars, char)
				}
			}
		}
	}
	return chars
}

func (f *Flds) SetMDTs(state bool) {
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		if ok {
			sf.Attrs.MDT = state
		}
	}
}

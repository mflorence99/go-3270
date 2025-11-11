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

// 👇 look at pairs of fields, this one and the next,
//
//	and assign all the cells in between to this field
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

// 🟦 Read functions

func (f *Flds) RB() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.flds {
		switch f.buf.Mode() {
		case types.FIELD_MODE:
			chars = append(chars, f.rbSF(fld)...)
		case types.EXTENDED_FIELD_MODE:
			chars = append(chars, f.rbSFE(fld)...)
		case types.CHARACTER_MODE:
			chars = append(chars, f.rbSA(fld)...)
		}
	}
	return chars
}

func (f *Flds) rbSF(fld Fld) []byte {
	chars := make([]byte, 0)
	sf, ok := fld.FldStart()
	if ok {
		chars = append(chars, byte(types.SF))
		chars = append(chars, sf.Attrs.Byte())
		for ix := 1; ix < len(fld); ix++ {
			cell := fld[ix]
			char := cell.Char
			if char != 0x00 {
				chars = append(chars, char)
			}
		}
	}
	return chars
}

func (f *Flds) rbSFE(fld Fld) []byte {
	chars := make([]byte, 0)
	sf, ok := fld.FldStart()
	if ok {
		chars = append(chars, byte(types.SFE))
		attrs := sf.Attrs.Bytes()
		chars = append(chars, byte(len(attrs)/2))
		chars = append(chars, attrs...)
		for ix := 1; ix < len(fld); ix++ {
			cell := fld[ix]
			char := cell.Char
			if char != 0x00 {
				chars = append(chars, char)
			}
		}
	}
	return chars
}

func (f *Flds) rbSA(fld Fld) []byte {
	chars := make([]byte, 0)
	_, ok := fld.FldStart()
	if ok {
		// TODO 🔥 support character mode
	}
	return chars
}

// TODO 🔥 CHARACTER_MODE *not* coded
func (f *Flds) RM() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		if ok && sf.Attrs.MDT {
			chars = append(chars, byte(types.SBA))
			chars = append(chars, conv.Addr2Bytes(sf.FldAddr+1)...)
			for ix := 1; ix < len(fld); ix++ {
				cell := fld[ix]
				char := cell.Char
				if char != 0x00 {
					chars = append(chars, char)
				}
			}
		}
	}
	return chars
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

func (f *Flds) ResetMDTs() {
	for _, fld := range f.flds {
		sf, ok := fld.FldStart()
		if ok {
			sf.Attrs.MDT = false
		}
	}
}

func (f *Flds) SetMDT(cell *Cell) bool {
	fld, ok := f.buf.Peek(cell.FldAddr)
	if !fld.FldStart || !ok {
		return false
	}
	fld.Attrs.MDT = true
	return true
}

package buffer

import (
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

// 🟦 Public functions

func (f *Flds) EAU() int {
	addr := -1
	for _, fld := range f.flds {
		sf, ok := fld.StartFld()
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

func (f *Flds) Get() []Fld {
	return f.flds
}

// TODO 🔥 *only* FIELD_MODE *not* coded
func (f *Flds) ReadBuffer() []byte {
	chars := make([]byte, 0)
	for _, fld := range f.flds {
		sf, _ := fld.StartFld()
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
		sf, _ := fld.StartFld()
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

func (f *Flds) Reset() {
	f.reset()
	// 👇 prepare to gather flds
	first := -1
	fld := make(Fld, 0)
	// 👇 normalize the cell at a given position in a given fld
	fix := func(fld Fld, cell *Cell) *Cell {
		sf, _ := fld.StartFld()
		if cell.FldAddr != sf.FldAddr {
			cell.Attrs = sf.Attrs
			cell.FldAddr = sf.FldAddr
		}
		return cell
	}
	// 👇 start at the beginning
	for ix := 0; ix < f.buf.Len(); ix++ {
		cell, _ := f.buf.Peek(ix)
		// 👇 a field is delimited by the next field
		if cell.FldStart {
			cell.FldAddr = ix
			if len(fld) > 0 {
				f.flds = append(f.flds, fld)
				fld = make(Fld, 0)
			}
			fld = append(fld, cell)
			// 👇 bookmark where we found the first field
			if first == -1 {
				first = ix
			}
		} else if first != -1 {
			cell = fix(fld, cell)
			fld = append(fld, cell)
		}
	}
	// 🔥 don't forget the last field, which includes any wrap-around
	if len(fld) > 0 {
		for ix := 0; ix < first; ix++ {
			cell, _ := f.buf.Peek(ix)
			cell = fix(fld, cell)
			fld = append(fld, cell)
		}
		f.flds = append(f.flds, fld)
	}
}

func (f *Flds) ResetMDTs() {
	for _, fld := range f.flds {
		sf, ok := fld.StartFld()
		if ok {
			sf.Attrs.Modified = false
		}
	}
}

func (f *Flds) SetMDT(fldAddr int) bool {
	fld, ok := f.buf.Peek(fldAddr)
	if !fld.FldStart || !ok {
		return false
	}
	fld.Attrs.Modified = true
	return true
}

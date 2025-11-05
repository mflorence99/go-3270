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

// 🟦 Housekeeping functions

func (f *Flds) Build() {
	// 👇 prepare to gather flds
	first := -1
	fld := make(Fld, 0)
	f.flds = make([]Fld, 0)
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
			fld = append(fld, cell)
		}
	}
	// 🔥 don't forget the last field, and include any cells found before the first SF or SFE
	if len(fld) > 0 {
		sf, _ := fld.FldStart()
		for ix := 0; ix < first; ix++ {
			cell, _ := f.buf.Peek(ix)
			cell.Attrs = sf.Attrs
			cell.FldAddr = sf.FldAddr
			fld = append(fld, cell)
		}
		f.flds = append(f.flds, fld)
	}
}

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

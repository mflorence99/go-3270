package buffer

import (
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
)

type Flds struct {
	buf  *Buffer
	bus  *pubsub.Bus
	cfg  pubsub.Config
	flds []Fld
}

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

// 🟦 Public methods

func (f *Flds) Get() []Fld {
	return f.flds
}

func (f *Flds) ReadMDT() []byte {
	bytes := make([]byte, 0)
	for _, fld := range f.flds {
		sf, _ := fld.StartFld()
		if sf.Attrs.Modified {
			bytes = append(bytes, byte(consts.SBA))
			bytes = append(bytes, conv.AddrToBytes(sf.FldAddr+1)...)
			for ix := 1; ix < len(fld); ix++ {
				cell := fld[ix]
				char := cell.Char
				if char == 0x00 {
					break
				}
				bytes = append(bytes, conv.A2E(char))
			}
		}
	}
	return bytes
}

func (f *Flds) Reset() {
	f.reset()
	// 👇 prepare to build flds
	stop := -1
	fld := make(Fld, 0)
	// 👇 start with an arbitrary cell
	cell, addr := f.buf.Get()
	start := addr
	for {
		// 👇 a field is delimited by the next field
		if cell != nil && cell.FldStart {
			if len(fld) > 0 {
				f.flds = append(f.flds, fld)
				fld = make(Fld, 0)
			}
			fld = append(fld, cell)
			// 👇 now we know where to stop
			if stop == -1 {
				stop = addr
			}
		} else if stop != -1 {
			sf, _ := fld.StartFld()
			// 👇 we are starting to collect cells now
			if cell == nil {
				// 👇 the cell might be missing, so inherit from the SF
				cell = &Cell{Attrs: sf.Attrs, Char: 0x00, FldAddr: sf.FldAddr}
				f.buf.Replace(cell, addr)
			} else if cell.FldAddr != sf.FldAddr {
				// 👇 the cell might be a residue of an earlier field after a W command
				cell.Attrs = sf.Attrs
				cell.FldAddr = sf.FldAddr
			}
			// 👇 finally, just append this cell to the fld
			fld = append(fld, cell)
		}
		// 👇 watch for wrap around as we blast through to stop
		cell, addr = f.buf.GetNext()
		if addr == stop {
			f.buf.Seek(start)
			break
		}
		f.buf.Seek(addr)
	}
	// 🔥 don't forget the last field!
	if len(fld) > 0 {
		f.flds = append(f.flds, fld)
	}
}

func (f *Flds) ResetMDT() {
	for _, fld := range f.flds {
		sf, ok := fld.StartFld()
		if ok {
			sf.Attrs.Modified = false
		}
	}
}

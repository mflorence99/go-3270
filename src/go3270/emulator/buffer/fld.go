package buffer

import (
	"strings"
)

// 🟧 Field in buffer

type Fld []*Cell

// 🟦 Public functions

func (f Fld) FldEnd() (*Cell, bool) {
	if len(f) > 0 {
		return f[len(f)-1], true
	}
	return nil, false
}

func (f Fld) FldStart() (*Cell, bool) {
	if len(f) > 0 {
		return f[0], true
	}
	return nil, false
}

// 🟦 Stringer implementation

func (f Fld) String() string {
	var b strings.Builder
	for ix := 1; ix < len(f); ix++ {
		cell := f[ix]
		if cell.Char >= ' ' {
			b.WriteByte(cell.Char)
		}
	}
	return strings.TrimSpace(b.String())
}

package buffer

import (
	"go3270/emulator/consts"
)

// 🟧 Cell in buffer

type Cell struct {
	Attrs    *consts.Attrs
	Char     byte
	FldAddr  int
	FldStart bool
	FldEnd   bool
}

// 🟦 Constructor

func NewCell() *Cell {
	c := new(Cell)
	c.Attrs = consts.DEFAULT_ATTRS
	return c
}

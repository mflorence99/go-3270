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
	FldGen   int
}

// 🟦 Constructor

func NewCell() *Cell {
	c := new(Cell)
	c.Attrs = &consts.Attrs{Default: true}
	return c
}

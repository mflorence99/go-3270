package buffer

import (
	"go3270/emulator/attrs"
)

// 🟧 Cell in buffer

type Cell struct {
	Attrs    *attrs.Attrs
	Char     byte
	FldAddr  int
	FldStart bool
	FldEnd   bool
}

// 🟦 Constructor

func NewCell() *Cell {
	c := new(Cell)
	c.Attrs = attrs.NewBasic(0)
	return c
}

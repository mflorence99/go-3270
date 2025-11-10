package buffer

import "go3270/emulator/types"

// 🟧 Cell in buffer

type Cell struct {
	Attrs    *types.Attrs
	Char     byte
	FldAddr  int
	FldStart bool
	FldEnd   bool
}

// 🟦 Constructor

func NewCell() *Cell {
	c := new(Cell)
	c.Attrs = types.DEFAULT_ATTRS
	return c
}

package buffer

import (
	"emulator/attrs"
)

type Cell struct {
	Attrs    *attrs.Attrs
	Char     byte
	FldAddr  int
	FldStart bool
}

func NewCell() *Cell {
	c := new(Cell)
	return c
}

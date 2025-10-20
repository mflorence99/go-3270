package attrs

type Cell struct {
	Attrs    *Attrs
	Char     byte
	FldAddr  int
	FldStart bool
}

func NewCell() *Cell {
	c := new(Cell)
	return c
}

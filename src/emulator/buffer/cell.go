package buffer

import (
	"emulator/attrs"
)

type Cell struct {
	Attrs *attrs.Attributes
	Char  byte
}

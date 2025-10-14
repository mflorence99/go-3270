package screen

import (
	"fmt"
)

type Box struct {
	X        float64
	Y        float64
	W        float64
	H        float64
	Baseline float64
}

func (b Box) String() string {
	return fmt.Sprintf("xyb[%f %f %f] wh[%f %f] ", b.X, b.Y, b.Baseline, b.W, b.H)
}

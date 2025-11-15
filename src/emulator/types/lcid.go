package types

import "fmt"

type LCID byte

// 🟦 Stringer implementation

func (l LCID) String() string {
	return fmt.Sprintf("%02x", byte(l))
}

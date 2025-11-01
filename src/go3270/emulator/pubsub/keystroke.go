package pubsub

import (
	"fmt"
	"go3270/emulator/utils"
)

// 🟧 Keystroke, as submitted by Typescript UI

type Keystroke struct {
	ALT   bool
	Code  string
	CTRL  bool
	Key   string
	SHIFT bool
}

// 🟦 Stringer implementation

func (k Keystroke) String() string {
	str := ""
	if k.CTRL {
		str += "CTRL "
	}
	if k.SHIFT {
		str += "SHIFT "
	}
	if k.ALT {
		str += "ALT "
	}
	return fmt.Sprintf("%s%s %s", str, k.Key, utils.Ternary(k.Code != k.Key && len(k.Key) > 1, k.Code, ""))
}

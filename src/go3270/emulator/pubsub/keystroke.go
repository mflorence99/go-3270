package pubsub

import (
	"fmt"
	"go3270/emulator/utils"
	"strings"
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
	var b strings.Builder
	if k.CTRL {
		b.WriteString("CTRL ")
	}
	if k.SHIFT {
		b.WriteString("SHIFT ")
	}
	if k.ALT {
		b.WriteString("ALT ")
	}
	return fmt.Sprintf("%s%s %s", b.String(), k.Key, utils.Ternary(k.Code != k.Key && len(k.Key) > 1, k.Code, ""))
}

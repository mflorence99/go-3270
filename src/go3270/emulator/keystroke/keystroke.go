package keystroke

import (
	"fmt"
	"go3270/emulator/bus"
	"go3270/emulator/utils"
)

type Keystroke struct {
	bus *bus.Bus
}

func NewKeystroke(bus *bus.Bus) *Keystroke {
	k := new(Keystroke)
	k.bus = bus
	k.bus.Subscribe("close", k.close)
	k.bus.Subscribe("keystroke", k.handler)
	return k
}

func (k *Keystroke) close() {}

func (k *Keystroke) handler(code string, key string, alt bool, ctrl bool, shift bool) {
	k.log(code, key, alt, ctrl, shift)
}

func (k *Keystroke) log(code string, key string, alt bool, ctrl bool, shift bool) {
	str := "⌨️ "
	if ctrl {
		str += "CTRL+"
	}
	if shift {
		str += "SHIFT+"
	}
	if alt {
		str += "ALT+"
	}
	println(fmt.Sprintf("%s%s %s", str, key, utils.Ternary(code != key && len(key) > 1, code, "")))
}

package keyboard

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"
)

type Keyboard struct {
	bus *pubsub.Bus
	cfg pubsub.Config
	buf *buffer.Buffer
	st  *state.State
}

func NewKeyboard(bus *pubsub.Bus, buf *buffer.Buffer, st *state.State) *Keyboard {
	k := new(Keyboard)
	k.bus = bus
	k.buf = buf
	k.st = st
	// 🔥 configure first
	k.bus.SubConfig(k.configure)
	k.bus.SubKeystroke(k.keystroke)
	k.bus.SubFocus(k.focus)
	return k
}

func (k *Keyboard) configure(cfg pubsub.Config) {
	k.cfg = cfg
}

func (k *Keyboard) focus(focussed bool) {
	println(fmt.Sprintf("⌨️ 3270 %s focus", utils.Ternary(focussed, "gains", "loses")))
	k.st.Patch(state.Patch{
		Error:   utils.BoolPtr(!focussed),
		Locked:  utils.BoolPtr(!focussed),
		Message: utils.StringPtr(utils.Ternary(focussed, "", "LOCK")),
	})
}

func (k *Keyboard) keystroke(key pubsub.Keystroke) {
	println(fmt.Sprintf("⌨️ %s", key))
}

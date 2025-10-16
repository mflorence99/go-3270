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
	// 👇 subscriptions
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
	cursorAt := k.st.Stat.CursorAt
	cursorTo := cursorAt
	cursorMax := k.cfg.Rows * k.cfg.Cols

	k.buf.Dirty.Push(cursorAt)

	switch key.Code {

	case "ArrowDown":
		cursorTo = cursorAt + k.cfg.Cols
		if cursorTo >= cursorMax {
			cursorTo = cursorAt % k.cfg.Cols
		}

	case "ArrowLeft":
		cursorTo = cursorAt - 1
		if cursorTo < 0 {
			cursorTo = cursorMax - 1
		}

	case "ArrowRight":
		cursorTo = cursorAt + 1
		if cursorTo >= cursorMax {
			cursorTo = 0
		}
	case "ArrowUp":
		cursorTo = cursorAt - k.cfg.Cols
		if cursorTo < 0 {
			cursorTo = (cursorAt % k.cfg.Cols) + cursorMax - k.cfg.Cols
		}

	}

	k.buf.Dirty.Push(cursorTo)

	k.st.Patch(state.Patch{
		CursorAt: utils.IntPtr(cursorTo),
	})
	k.buf.Seek(cursorTo)

	k.bus.PubRender()
}

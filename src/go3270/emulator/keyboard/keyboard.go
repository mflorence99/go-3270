package keyboard

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
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
	k.buf = buf
	k.bus = bus
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
	// 👇 prepare to move the cursor -- most keystrokes do this
	cursorAt := k.st.Stat.CursorAt
	cursorTo := cursorAt
	cursorMax := k.cfg.Rows * k.cfg.Cols
	// 👇 maintain a stack of changed cells
	deltas := utils.NewStack[int](1)
	// 👇 make sure we know where to start
	k.buf.Seek(cursorAt)
	// 👇 pre-analyze AID key
	aid := consts.AIDOf(key.Key, key.ALT, key.CTRL, key.SHIFT)

	switch {

	case aid == consts.CLEAR:
		println(fmt.Sprintf("🐞 %s", aid))
		k.bus.PubInboundAttn(aid)

	case aid == consts.ENTER:
		println(fmt.Sprintf("🐞 %s", aid))
		k.bus.PubInboundRM(aid)

	case aid.PAx():
		println(fmt.Sprintf("🐞 %s", aid))
		k.bus.PubInboundAttn(aid)

	case aid.PFx():
		println(fmt.Sprintf("🐞 %s", aid))
		k.bus.PubInboundRM(aid)

	case key.Code == "ArrowDown":
		cursorTo = cursorAt + k.cfg.Cols
		if cursorTo >= cursorMax {
			cursorTo = cursorAt % k.cfg.Cols
		}

	case key.Code == "ArrowLeft":
		cursorTo = cursorAt - 1
		if cursorTo < 0 {
			cursorTo = cursorMax - 1
		}

	case key.Code == "ArrowRight":
		cursorTo = cursorAt + 1
		if cursorTo >= cursorMax {
			cursorTo = 0
		}
	case key.Code == "ArrowUp":
		cursorTo = cursorAt - k.cfg.Cols
		if cursorTo < 0 {
			cursorTo = (cursorAt % k.cfg.Cols) + cursorMax - k.cfg.Cols
		}

	case key.Code == "Backspace":
		_, ok := k.buf.Backspace()
		if ok {
			cursorTo = k.buf.Addr()
		} else {
			k.st.Patch(state.Patch{Alarm: utils.BoolPtr(true)})
		}

	case key.Code == "Tab":
		println(fmt.Sprintf("🐞 tab %s", utils.Ternary(key.SHIFT, "bwd", "fwd")))
		_, ok := k.buf.Tab(utils.Ternary(key.SHIFT, -1, +1))
		if ok {
			cursorTo = k.buf.Addr()
		} else {
			k.st.Patch(state.Patch{Alarm: utils.BoolPtr(true)})
		}

	case len(key.Key) == 1:
		_, ok := k.buf.Keyin(key.Key[0])
		if ok {
			cursorTo = k.buf.Addr()
		} else {
			k.st.Patch(state.Patch{Alarm: utils.BoolPtr(true)})
		}
	}

	// 👇 only if the cursor has moved!
	if cursorTo != cursorAt {
		deltas.Push(cursorAt)
		deltas.Push(cursorTo)
		k.buf.Seek(cursorTo)
		// 👇 update the status depending on the new cell
		cell, _ := k.buf.Get()
		k.st.Patch(state.Patch{
			CursorAt:  utils.IntPtr(cursorTo),
			Numeric:   utils.BoolPtr(cell.Attrs.Numeric),
			Protected: utils.BoolPtr(cell.Attrs.Protected || cell.FldStart),
		})
	}
	// 👇 render any changes
	if !deltas.Empty() {
		k.bus.PubRender(deltas)
	}
}

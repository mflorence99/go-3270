package keyboard

import (
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"

	"strings"
)

type Keyboard struct {
	bus  *pubsub.Bus
	buf  *buffer.Buffer
	cfg  pubsub.Config
	flds *buffer.Flds
	st   *state.State
}

func NewKeyboard(bus *pubsub.Bus, buf *buffer.Buffer, flds *buffer.Flds, st *state.State) *Keyboard {
	k := new(Keyboard)
	k.buf = buf
	k.bus = bus
	k.flds = flds
	k.st = st
	// 👇 subscriptions
	k.bus.SubConfig(k.configure)
	k.bus.SubKeystroke(k.keystroke)
	k.bus.SubFocus(k.focus)
	return k
}

func (k *Keyboard) backspace() (int, bool) {
	// 👇 validate data entry into previous cell
	c, addr := k.buf.PrevGet()
	prot := c.FldStart || c.Attrs.Protected
	if prot {
		return -1, false
	}
	// 👇 reposition to previous cell and update it
	k.buf.Seek(addr)
	c.Char = ' '
	c.Attrs.Modified = true
	addr = k.buf.Set(c)
	// 👇 set the MDT flag at the field level
	ok := k.flds.SetMDT(c.FldAddr)
	if !ok {
		return -1, false
	}
	return addr, true
}

func (k *Keyboard) configure(cfg pubsub.Config) {
	k.cfg = cfg
}

func (k *Keyboard) focus(focussed bool) {
	k.st.Patch(state.Patch{
		Error:   utils.BoolPtr(!focussed),
		Locked:  utils.BoolPtr(!focussed),
		Message: utils.StringPtr(utils.Ternary(focussed, "", "LOCK")),
	})
}

func (k *Keyboard) keyin(char byte) (int, bool) {
	c, _ := k.buf.Get()
	// 👇 validate data entry into current cell
	numlock := c.Attrs.Numeric && !strings.Contains("-0123456789.", string(char))
	prot := c.FldStart || c.Attrs.Protected
	if numlock || prot {
		return -1, false
	}
	// 👇 update cell and advance to next
	c.Char = char
	c.Attrs.Modified = true
	addr := k.buf.SetAndNext(c)
	// 👇 set the MDT flag at the field level
	ok := k.flds.SetMDT(c.FldAddr)
	if !ok {
		return -1, false
	}
	return addr, true
}

func (k *Keyboard) keystroke(key pubsub.Keystroke) {
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
		k.bus.PubAttn(aid)

	case aid == consts.ENTER:
		k.bus.PubRM(aid)

	case aid.PAx():
		k.bus.PubAttn(aid)

	case aid.PFx():
		k.bus.PubRM(aid)

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
		_, ok := k.backspace()
		if ok {
			cursorTo = k.buf.Addr()
		} else {
			k.st.Patch(state.Patch{Alarm: utils.BoolPtr(true)})
		}

	case key.Code == "Tab":
		_, ok := k.tab(utils.Ternary(key.SHIFT, -1, +1))
		if ok {
			cursorTo = k.buf.Addr()
		} else {
			k.st.Patch(state.Patch{Alarm: utils.BoolPtr(true)})
		}

	case len(key.Key) == 1:
		_, ok := k.keyin(key.Key[0])
		if ok {
			cursorTo = k.buf.Addr()
		} else {
			k.st.Patch(state.Patch{Alarm: utils.BoolPtr(true)})
		}
	}

	// 👇 probe cursor position for debugging
	if key.CTRL && strings.HasPrefix(key.Code, "Arrow") {
		k.bus.PubProbe(cursorTo)
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
		k.bus.PubRenderDeltas(deltas)
	}
}

func (k *Keyboard) tab(dir int) (int, bool) {
	start := k.buf.Addr()
	addr := k.buf.Addr()
	for ix := 0; ; ix++ {
		// 👇 wrap to the start means no unprotected field
		if addr == start && ix > 0 {
			break
		}
		// 👇 look at the "next" cell
		addr += dir
		if addr < 0 {
			addr = k.buf.Len() - 1
		} else if addr >= k.buf.Len() {
			addr = 0
		}
		// 👇 see if we've hit an unprotected field start
		cell, _ := k.buf.Peek(addr)
		if cell.FldStart && !cell.Attrs.Protected {
			// 👇 if going backwards, and we hit in the first try, it doesn't count
			if dir < 0 && ix == 0 {
				continue
			}
			k.buf.Seek(addr) // 👈 go to FldStart
			cell, addr := k.buf.GetNext()
			// 👇 if the next cell is also a field start (two contiguous SFs) it also doesn't count
			if cell.FldStart {
				continue
			}
			k.buf.Seek(addr) // 👈 now to first char
			return addr, true
		}
	}
	return -1, false
}

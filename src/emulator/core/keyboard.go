package core

import (
	"emulator/conv"
	"emulator/types"
	"emulator/utils"

	"strings"
)

// 🟧 Respond to keyboard input

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

// 👁️ Keyboard Operations pp 7-10 to 7-15

type Keyboard struct {
	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewKeyboard(emu *Emulator) *Keyboard {
	k := new(Keyboard)
	k.emu = emu
	// 👇 subscriptions
	k.emu.Bus.SubInit(k.init)
	k.emu.Bus.SubKeystroke(k.keystroke)
	k.emu.Bus.SubFocus(k.focus)
	return k
}

// TODO 🔥 just in case we need it
func (k *Keyboard) init() {}

// 🟦 Gain/lose focus

func (k *Keyboard) focus(focussed bool) {
	k.emu.State.Patch(types.Patch{
		Error:   utils.BoolPtr(!focussed),
		Locked:  utils.BoolPtr(!focussed),
		Message: utils.StringPtr(utils.Ternary(focussed, "", "LOCK")),
	})
}

// 🟦 Dispatch action per key code

func (k *Keyboard) keystroke(key types.Keystroke) {
	// 👇 prepare to move the cursor -- most keystrokes do this
	cursorAt := k.emu.State.Status.CursorAt
	cursorTo := cursorAt
	cursorMax := k.emu.Cfg.Rows * k.emu.Cfg.Cols
	// 👇 maintain a stack of changed cells
	deltas := utils.NewStack[uint](1)
	// 👇 make sure we know where to start
	k.emu.Buf.WrappingSeek(int(cursorAt))
	// 👇 pre-analyze AID key
	aid := types.AIDOf(key.Key, key.ALT, key.CTRL, key.SHIFT)

	switch {

	case aid == types.CLEAR:
		// TODO 🔥 this is a backdoor for testing
		if key.SHIFT {
			k.emu.Bus.PubRB(aid)
		} else {
			k.emu.Bus.PubReset()
			k.emu.Bus.PubAttn(aid)
		}

	case aid == types.ENTER:
		k.emu.Bus.PubRM(aid)

	case aid.PAx():
		k.emu.Bus.PubAttn(aid)

	case aid.PFx():
		k.emu.Bus.PubRM(aid)

	case key.Code == "ArrowDown":
		if cursorAt >= cursorMax-k.emu.Cfg.Cols {
			cursorTo = cursorAt % k.emu.Cfg.Cols
		} else {
			cursorTo = cursorAt + k.emu.Cfg.Cols
		}

	case key.Code == "ArrowLeft":
		if cursorAt == 0 {
			cursorTo = cursorMax - 1
		} else {
			cursorTo = cursorAt - 1
		}

	case key.Code == "ArrowRight":
		if cursorAt == cursorMax-1 {
			cursorTo = 0
		} else {
			cursorTo = cursorAt + 1
		}

	case key.Code == "ArrowUp":
		if cursorAt < k.emu.Cfg.Cols {
			cursorTo = (cursorAt % k.emu.Cfg.Cols) + cursorMax - k.emu.Cfg.Cols
		} else {
			cursorTo = cursorAt - k.emu.Cfg.Cols
		}

	case key.Code == "Backspace":
		addr, ok := k.backspace()
		if ok {
			k.emu.Buf.MustSeek(addr)
			cursorTo = addr
		} else {
			k.emu.State.Patch(types.Patch{Alarm: utils.BoolPtr(true)})
		}

	case key.Code == "Tab":
		_, ok := k.tab(utils.Ternary(key.SHIFT, -1, +1))
		if ok {
			cursorTo = k.emu.Buf.Addr()
		} else {
			k.emu.State.Patch(types.Patch{Alarm: utils.BoolPtr(true)})
		}

	case len(key.Key) == 1:
		addr, ok := k.keyin(key.Key[0])
		if ok {
			cursorTo = addr
		} else {
			k.emu.State.Patch(types.Patch{Alarm: utils.BoolPtr(true)})
		}
	}

	// 👇 probe cursor position for debugging
	if key.CTRL && strings.HasPrefix(key.Code, "Arrow") {
		k.emu.Bus.PubProbe(cursorTo)
	}

	// 👇 only if the cursor has moved!
	if cursorTo != cursorAt {
		deltas.Push(cursorAt)
		deltas.Push(cursorTo)
		k.emu.Buf.MustSeek(cursorTo)
		// 👇 update the status depending on the new cell
		cell, _ := k.emu.Buf.Get()
		k.emu.State.Patch(types.Patch{
			CursorAt:  utils.UintPtr(cursorTo),
			Numeric:   utils.BoolPtr(cell.Attrs.Numeric),
			Protected: utils.BoolPtr(cell.Attrs.Protected || cell.IsFldStart()),
		})
	}
	// 👇 render any changes
	if !deltas.Empty() {
		k.emu.Bus.PubRenderDeltas(deltas)
	}
}

// 🟦 BACKSPACE

func (k *Keyboard) backspace() (uint, bool) {
	cell, _ := k.emu.Buf.Get()
	// 👇 validate data entry into current cell
	prot := cell.IsFldStart() || cell.Attrs.Protected
	if prot {
		return 0, false
	}
	// 👇 update cell
	cell.Char = 0x40
	cell.Attrs.MDT = true
	// 👇 set the MDT flag at the field level
	sf, ok := cell.GetFldStart()
	if !ok {
		return 0, false
	}
	sf.Attrs.MDT = true
	// 👇 if the previous cell is a field start, don't advance
	next, addr := k.emu.Buf.PrevGet()
	if next.IsFldStart() {
		return k.emu.Buf.Addr(), false
	}
	// 👇 advance to previous cell
	return k.emu.Buf.MustSeek(addr), true
}

// 🟦 KEYSTROKE

func (k *Keyboard) keyin(char byte) (uint, bool) {
	cell, _ := k.emu.Buf.Get()
	// 👇 validate data entry into current cell
	numlock := cell.Attrs.Numeric && !strings.Contains("-0123456789.", string(char))
	prot := cell.IsFldStart() || cell.Attrs.Protected
	if numlock || prot {
		return 0, false
	}
	// 👇 update cell and advance to next
	cell.Char = conv.A2E(char)
	cell.Attrs.MDT = true
	// 👇 set the MDT flag at the field level
	sf, ok := cell.GetFldStart()
	if !ok {
		return 0, false
	}
	sf.Attrs.MDT = true
	// 👇 if the next cell is a field start with autoskip, tab to next Fld
	next, addr := k.emu.Buf.GetNext()
	if next.IsFldStart() {
		if next.Attrs.Autoskip {
			return k.tab(+1)
		}
		// 👇 don't advance if not autoskip
		return k.emu.Buf.Addr(), false
	}
	// 👇 advance to next cell
	return k.emu.Buf.MustSeek(addr), true
}

// 🟦 TAB

func (k *Keyboard) tab(dir int) (uint, bool) {
	advance := func(dir int) (*Cell, uint) {
		if dir > 0 {
			return k.emu.Buf.GetNext()
		}
		return k.emu.Buf.PrevGet()
	}
	// 🔥 look in opposite direction for the stop addr
	start := k.emu.Buf.Addr()
	cell, stop := advance(dir * -1)
	// 🔥 we don't really need to seek here, but we're following a pattern
	addr := k.emu.Buf.MustSeek(start)
	// 👇 keep looking for an unprotected Fld start until we wrap
	for addr != stop {
		cell, addr = advance(dir)
		if cell.IsFldStart() && !cell.Attrs.Protected {
			return k.emu.Buf.WrappingSeek(int(addr) + 1), true
		}
		k.emu.Buf.MustSeek(addr)
	}
	// 👇 we wrapped all the way around to the start w/o unprotected
	return k.emu.Buf.MustSeek(start), false
}

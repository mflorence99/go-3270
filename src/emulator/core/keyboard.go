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

// 🟦 Dispatch action per key code

func (k *Keyboard) keystroke(key types.Keystroke) {
	// 👇 prepare to move the cursor -- most keystrokes do this
	cursorAt := k.emu.State.Status.CursorAt
	cursorTo := cursorAt
	cursorMax := k.emu.Cfg.Rows * k.emu.Cfg.Cols
	// 👇 maintain a stack of changed cells
	deltas := utils.NewStack[uint](1)
	// 👇 make sure we know where to start
	k.emu.Buf.WrappingSeek(cursorAt)
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
		cursorTo = cursorAt + k.emu.Cfg.Cols
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
		cursorTo = cursorAt + 1
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
		_, ok := k.backspace()
		if ok {
			cursorTo = k.emu.Buf.Addr()
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
		_, ok := k.keyin(key.Key[0])
		if ok {
			cursorTo = k.emu.Buf.Addr()
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

// 🟦 Functions specific to particular keys

func (k *Keyboard) backspace() (uint, bool) {
	// 👇 validate data entry into previous cell
	cell, addr := k.emu.Buf.PrevGet()
	prot := cell.IsFldStart() || cell.Attrs.Protected
	if prot {
		return 0, false
	}
	// 👇 reposition to previous cell and update it
	k.emu.Buf.MustSeek(addr)
	cell.Char = 0x40
	cell.Attrs.MDT = true
	// 👇 set the MDT flag at the field level
	sf, ok := cell.GetFldStart()
	if !ok {
		return 0, false
	}
	sf.Attrs.MDT = true
	return addr, true
}

func (k *Keyboard) focus(focussed bool) {
	k.emu.State.Patch(types.Patch{
		Error:   utils.BoolPtr(!focussed),
		Locked:  utils.BoolPtr(!focussed),
		Message: utils.StringPtr(utils.Ternary(focussed, "", "LOCK")),
	})
}

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
	addr := k.emu.Buf.SetAndNext(cell)
	// 👇 set the MDT flag at the field level
	// sf, ok := k.emu.Buf.Peek(cell.FldAddr)
	// if !ok {
	// 	return 0, false
	// }
	// sf.Attrs.MDT = true
	// 👇 if we typed into a field start with autoskip, tab to next
	next, _ := k.emu.Buf.Get()
	if next.IsFldStart() && next.Attrs.Autoskip {
		return k.tab(+1)
	}
	return addr, true
}

func (k *Keyboard) tab(dir int) (uint, bool) {
	start := k.emu.Buf.Addr()
	addr := k.emu.Buf.Addr()
	for ix := 0; ; ix++ {
		// 👇 wrap to the start means no unprotected field
		if addr == start && ix > 0 {
			break
		}
		// 👇 look at the "next" cell
		// addr += dir
		if addr < 0 {
			addr = k.emu.Buf.Len() - 1
		} else if addr >= k.emu.Buf.Len() {
			addr = 0
		}
		// 👇 see if we've hit an unprotected field start
		cell := k.emu.Buf.MustPeek(addr)
		if cell.IsFldStart() && !cell.Attrs.Protected {
			// 👇 if going backwards, and we hit in the first try, it doesn't count
			if dir < 0 && ix == 0 {
				continue
			}
			k.emu.Buf.WrappingSeek(addr) // 👈 go to FldStart
			cell, addr := k.emu.Buf.GetNext()
			// 👇 if the next cell is also a field start (two contiguous SFs)
			//    it also doesn't count
			if cell.IsFldStart() {
				continue
			}
			k.emu.Buf.WrappingSeek(addr) // 👈 now to first char
			return addr, true
		}
	}
	return 0, false
}

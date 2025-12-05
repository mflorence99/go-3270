package core

import (
	"emulator/conv"
	"emulator/types"
	"emulator/utils"

	"slices"
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
	k.emu.Bus.SubKeystroke(k.keystroke)
	k.emu.Bus.SubFocus(k.focus)
	return k
}

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
	// 👇 see if we are in insert mode
	insertMode := k.emu.State.Status.Insert
	// 👇 prepare to move the cursor -- many keystrokes do this
	cursorAt := k.emu.State.Status.CursorAt
	cursorTo := cursorAt
	cursorMax := k.emu.Cfg.Rows * k.emu.Cfg.Cols
	// 👇 maintain a stack of changed cells
	deltas := utils.NewStack[uint](1)
	// 👇 make sure we know where to start
	k.emu.Buf.WrappingSeek(int(cursorAt))
	// 👇 pre-analyze AID key
	aid := types.AIDOf(key.Key, key.ALT, key.CTRL, key.SHIFT)
	// 👇 assume success of operation
	ok := true

	switch {

	case aid == types.CLEAR:
		// 🔥 this is a backdoor for testing
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
		cursorTo, ok = k.backspace(cursorAt)

	case key.Code == "Delete":
		cursorTo, ok = k.delete(cursorAt, deltas)

	case key.Code == "End":
		cursorTo, ok = k.end(cursorAt)

	case key.Code == "Home":
		cursorTo, ok = k.home(cursorAt)

	case key.Code == "Insert":
		k.emu.State.Patch(types.Patch{
			Insert: utils.BoolPtr(!insertMode),
		})

	case key.Code == "Tab":
		cursorTo, ok = k.tab(utils.Ternary(key.SHIFT, -1, +1), cursorAt)

	case len(key.Key) == 1:
		cursorTo, ok = k.keyin(key.Key[0], cursorAt, deltas, insertMode)

	}

	// 👇 something might've gone wrong
	if !ok {
		k.emu.State.Patch(types.Patch{Alarm: utils.BoolPtr(true)})
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

func (k *Keyboard) backspace(dfltAddr uint) (uint, bool) {
	cell, addr := k.emu.Buf.PrevGet()
	// 👇 validate data entry into previous cell
	prot := cell.IsFldStart() || cell.Attrs.Protected
	if prot {
		return dfltAddr, false
	}
	// 👇 set the MDT flag at the field level
	sf, ok := cell.GetFldStart()
	if !ok {
		return dfltAddr, false
	}
	sf.Attrs.MDT = true
	// 👇 update cell
	cell.Char = 0x40
	// 👇 advance to previous cell
	return k.emu.Buf.MustSeek(addr), true
}

// 🟦 DELETE

func (k *Keyboard) delete(dfltAddr uint, deltas *utils.Stack[uint]) (uint, bool) {
	cell, _ := k.emu.Buf.Get()
	// 👇 only if in unprotected field with at least two character cells
	fld, ok := cell.FindFld()
	if !ok || len(fld.Cells) <= 2 || fld.Cells[0].Attrs.Protected {
		return dfltAddr, false
	}
	// 👇 shift all subsequent characters from the right
	ix := 0
	iy := slices.Index(fld.Cells, cell)
	for ix = iy; ix < len(fld.Cells)-1; ix++ {
		fld.Cells[ix].Char = fld.Cells[ix+1].Char
	}
	// 👇 fill the remainder with nulls
	for ; ix < len(fld.Cells); ix++ {
		fld.Cells[ix].Char = 0x00
	}
	// 👇 set the MDT flag at the field level
	sf := fld.Cells[0]
	sf.Attrs.MDT = true
	// 👇 indicate ALL the cells that changed
	addr, _ := sf.GetFldAddr()
	for ix = iy; ix < len(fld.Cells); ix++ {
		deltas.Push(k.emu.Buf.WrapAddr(int(addr) + ix))
	}
	// 🔥 the cursor doesn't move in this operation
	return dfltAddr, true
}

// 🟦 END

func (k *Keyboard) end(dfltAddr uint) (uint, bool) {
	cell, _ := k.emu.Buf.Get()
	// 👇 only if in unprotected field with home cell
	fld, ok := cell.FindFld()
	if !ok || len(fld.Cells) <= 1 || fld.Cells[0].Attrs.Protected {
		return dfltAddr, false
	}
	// 👇 look backward for first non-blank, then position + 1 from it
	for ix, cell := range slices.Backward(fld.Cells) {
		if cell.Char > 0x40 {
			addr, _ := cell.GetFldAddr()
			eof := min(len(fld.Cells)-1, ix+1)
			return k.emu.Buf.WrappingSeek(int(addr) + eof), true
		}
	}
	// 🔥 not an error if whole field is blank
	return dfltAddr, true
}

// 🟦 HOME

func (k *Keyboard) home(dfltAddr uint) (uint, bool) {
	for _, fld := range k.emu.Flds.Flds {
		sf := fld.Cells[0]
		// 👇 find the first unprotected field with a home cell
		if !sf.Attrs.Protected && len(fld.Cells) > 1 {
			addr, _ := sf.GetFldAddr()
			return k.emu.Buf.WrappingSeek(int(addr) + 1), true
		}
	}
	// 👇 looked everywhere and couldn't find it!
	return dfltAddr, false
}

// 🟦 KEYSTROKE

func (k *Keyboard) keyin(char byte, dfltAddr uint, deltas *utils.Stack[uint], insertMode bool) (uint, bool) {
	cell, _ := k.emu.Buf.Get()
	if k.keyinvalid(cell, char) || !k.keyinMDT(cell) {
		return dfltAddr, false
	}
	if insertMode {
		return k.keyinsert(cell, char, dfltAddr, deltas)
	}
	return k.keyinover(cell, char, dfltAddr)
}

func (k *Keyboard) keyinvalid(cell *Cell, char byte) bool {
	numlock := cell.Attrs.Numeric && !strings.Contains("-0123456789.", string(char))
	prot := cell.IsFldStart() || cell.Attrs.Protected
	if numlock || prot {
		return true
	}
	return false
}

func (k *Keyboard) keyinMDT(cell *Cell) bool {
	sf, ok := cell.GetFldStart()
	if !ok {
		return false
	}
	sf.Attrs.MDT = true
	return true
}

func (k *Keyboard) keyinsert(cell *Cell, char byte, dfltAddr uint, deltas *utils.Stack[uint]) (uint, bool) {
	// 👇 can't insert if not in a field or if the field is full
	fld, ok := cell.FindFld()
	if !ok || fld.Cells[len(fld.Cells)-1].Char > 0x40 {
		return dfltAddr, false
	}
	// 👇 shift all subsequent characters to the right
	ix := 0
	iy := slices.Index(fld.Cells, cell)
	for ix = len(fld.Cells) - 1; ix > iy; ix-- {
		fld.Cells[ix].Char = fld.Cells[ix-1].Char
	}
	cell.Char = conv.A2E(char)
	// 👇 indicate ALL the cells that changed
	addr, _ := cell.GetFldAddr()
	for ix = iy; ix < len(fld.Cells); ix++ {
		deltas.Push(k.emu.Buf.WrapAddr(int(addr) + ix))
	}
	// 🔥 advance to next cell
	return k.emu.Buf.WrappingSeek(int(addr) + iy + 1), true
}

func (k *Keyboard) keyinover(cell *Cell, char byte, dfltAddr uint) (uint, bool) {
	cell.Char = conv.A2E(char)
	// 👇 if the next cell is a field start with autoskip, tab to next Fld
	next, addr := k.emu.Buf.GetNext()
	if next.IsFldStart() {
		if next.Attrs.Autoskip {
			return k.tab(+1, dfltAddr)
		}
		// 👇 don't advance if not autoskip
		return dfltAddr, false
	}
	// 👇 advance to next cell
	return k.emu.Buf.MustSeek(addr), true
}

// 🟦 TAB

func (k *Keyboard) tab(dir int, start uint) (uint, bool) {
	var cell *Cell
	advance := func(dir int) (*Cell, uint) {
		if dir > 0 {
			return k.emu.Buf.GetNext()
		}
		return k.emu.Buf.PrevGet()
	}
	// 🔥 look in opposite direction for the stop addr
	_, stop := advance(dir * -1)
	// 🔥 we don't really need to seek here, but we're following a pattern
	addr := k.emu.Buf.MustSeek(start)
	// 👇 keep looking for an unprotected Fld start until we wrap
	for addr != stop {
		cell, addr = advance(dir)
		if cell.IsFldHome() && !cell.Attrs.Protected {
			return k.emu.Buf.MustSeek(addr), true
		}
		k.emu.Buf.MustSeek(addr)
	}
	// 👇 we wrapped all the way around to the start w/o unprotected
	return k.emu.Buf.MustSeek(start), false
}

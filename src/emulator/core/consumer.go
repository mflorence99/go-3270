package core

import (
	"emulator/conv"
	"emulator/types"
	"emulator/utils"
	"fmt"
	"time"
)

// 🟧 Consume outbound (3270 <- app) data stream

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

type Consumer struct {
	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewConsumer(emu *Emulator) *Consumer {
	c := new(Consumer)
	c.emu = emu
	// 👇 subscriptions
	c.emu.Bus.SubOutbound(c.consume)
	return c
}

func (c *Consumer) consume(chars []byte) {
	defer utils.ElapsedTime(time.Now())
	// 👇 process the commands in the stream
	out := NewOutbound(chars, c.emu.Bus)
	char := out.MustNext()
	cmd := types.Command(char)
	c.commands(out, cmd)
	// 👇 once stream is processed we are able to reflect current cell status
	cursorAt := c.emu.State.Status.CursorAt
	cell := c.emu.Buf.MustPeek(cursorAt)
	c.emu.State.Patch(types.Patch{
		Insert:    utils.BoolPtr(false),
		Numeric:   utils.BoolPtr(cell.Attrs.Numeric),
		Protected: utils.BoolPtr(cell.Attrs.Protected || cell.IsFldStart()),
	})
}

// 🟦 Commands

func (c *Consumer) commands(out *Outbound, cmd types.Command) {
	// 👇 dispatch on command
	switch cmd {

	case types.EAU:
		c.eau()

	case types.EW:
		c.ew(out)

	case types.EWA:
		c.ewa(out)

	case types.RB:
		c.rb()

	case types.RM:
		c.rm()

	case types.RMA:
		c.rma()

	case types.W:
		c.w(out)

	case types.WSF:
		c.wsf(out)
	}
}

func (c *Consumer) eau() {
	if addr, ok := c.emu.Flds.EAU(); ok {
		c.emu.Buf.WrappingSeek(int(addr) + 1)
		c.emu.State.Patch(types.Patch{
			CursorAt: utils.UintPtr(c.emu.Buf.Addr()),
		})
	}
}

func (c *Consumer) ew(out *Outbound) {
	if c.wcc(out) {
		c.emu.Bus.PubReset()
		c.orders(out)
		c.emu.Bus.PubRender()
	}
}

func (c *Consumer) ewa(out *Outbound) {
	if c.wcc(out) {
		c.emu.Bus.PubReset()
		c.orders(out)
		c.emu.Bus.PubRender()
	}
}

func (c *Consumer) rb() {
	c.emu.Bus.PubRB(types.INBOUND)
}

func (c *Consumer) rm() {
	c.emu.Bus.PubRM(types.INBOUND)
}

func (c *Consumer) rma() {
	c.emu.Bus.PubRMA(types.INBOUND)
}

func (c *Consumer) w(out *Outbound) {
	c.wcc(out)
	c.orders(out)
	c.emu.Bus.PubRender()
}

func (c *Consumer) wcc(out *Outbound) bool {
	if char, ok := out.Next(); ok {
		wcc := types.NewWCC(char)
		// TODO 🔥 not yet handled
		if wcc.Reset {
			println("🔥 WCC Reset not implemented")
		}
		if wcc.ResetMDT {
			for _, fld := range c.emu.Flds.Flds {
				sf := fld.Cells[0]
				sf.Attrs.MDT = false
			}
		}
		c.emu.Bus.PubWCChar(wcc)
		return true
	} else {
		return false
	}
}

// 🟦 WSF (which may contain multiple commands itself)

func (c *Consumer) wsf(out *Outbound) {
	// TODO 🔥 there are a million SF types
	sflds := SFldsFromStream(out)
	for _, sfld := range sflds {

		switch sfld.ID {

		case types.READ_PARTITION:
			c.rp(sfld)

		case types.SET_REPLY_MODE:
			c.srm(sfld)

		default:
			c.emu.Bus.PubPanic(fmt.Sprintf("🔥 SFld %s not implemented", sfld))

		}
	}
}

func (c *Consumer) rp(sfld SFld) {
	pid := sfld.Info[0]
	if pid == 0xfF {
		cmd := sfld.Info[1]

		switch types.Command(cmd) {

		case types.Q:
			c.emu.Bus.PubQ()

		case types.QL:
			all := (sfld.Info[2] & 0b10000000) == 0b10000000
			var qcodes []types.QCode
			if all {
				qcodes = []types.QCode{
					types.USABLE_AREA,
					types.ALPHANUMERIC_PARTITIONS,
					types.CHARACTER_SETS,
					types.COLOR_SUPPORT,
					types.HIGHLIGHTING,
					types.REPLY_MODES,
					types.FIELD_VALIDATION,
					types.FIELD_OUTLINING,
					types.DDM,
					types.RPQ_NAMES,
					types.IMPLICIT_PARTITION,
				}
			} else {
				qcodes = make([]types.QCode, 0)
				for ix := 3; ix < len(sfld.Info); ix++ {
					qcodes = append(qcodes, types.QCode(sfld.Info[ix]))
				}
			}
			c.emu.Bus.PubQL(qcodes)

		case types.RB:
			c.emu.Bus.PubRB(types.INBOUND)

		case types.RM:
			c.emu.Bus.PubRM(types.INBOUND)

		case types.RMA:
			c.emu.Bus.PubRMA(types.INBOUND)

		}
	}
}

func (c *Consumer) srm(sfld SFld) {
	pid := sfld.Info[0]
	if pid == 0xfF {
		mode := types.Mode(sfld.Info[1])
		c.emu.Buf.SetMode(mode)
	}
}

// 🟦 Orders

func (c *Consumer) orders(out *Outbound) {
	var inFld bool
	var fldAddr uint
	fldAttrs := types.DEFAULT_ATTRS

	// 👇 look at each byte to see if it is an order
	for out.HasNext() {
		char := out.MustNext()
		order := types.Order(char)
		// 👇 dispatch on order
		switch order {

		case types.EUA:
			c.eua(out)

		case types.GE:
			c.ge(out, fldAddr, fldAttrs, inFld)

		case types.IC:
			c.ic()

		case types.MF:
			c.mf(out)

		case types.PT:
			c.pt()

		case types.RA:
			c.ra(out, fldAddr, fldAttrs, inFld)

		case types.SA:
			fldAttrs = c.sa(out, fldAttrs)

		case types.SBA:
			c.sba(out)

		case types.SF:
			inFld = true
			fldAddr, fldAttrs = c.sf(out)

		case types.SFE:
			inFld = true
			fldAddr, fldAttrs = c.sfe(out)

		// 👇 if it isn't an order, it's data
		default:
			c.char(char, fldAddr, fldAttrs, inFld)
		}
	}
}

func (c *Consumer) char(char byte, fldAddr uint, fldAttrs *types.Attrs, inFld bool) {
	cell := NewCell(c.emu)
	cell.Attrs = fldAttrs
	cell.Char = char
	if inFld {
		cell.SetFldAddr(fldAddr)
	}
	c.emu.Buf.SetAndNext(cell)
}

func (c *Consumer) eua(out *Outbound) {
	raw := out.MustNextSlice(2)
	stop := conv.Bytes2Addr(raw)
	// 👇 validate stop addr
	_, ok := c.emu.Buf.Peek(stop)
	if ok {
		c.emu.Cells.EUA(c.emu.Buf.Addr(), stop)
	} else {
		c.emu.Buf.AddrPanic(stop)
	}
}

func (c *Consumer) ge(out *Outbound, fldAddr uint, fldAttrs *types.Attrs, inFld bool) {
	char := out.MustNext()
	// TODO 🔥 GE not properly handled -- what alt character set??
	// also needs to be present in inbound stream (RB, RM/A)
	fldAttrs.LCID = 0xf1
	c.char(char, fldAddr, fldAttrs, inFld)
}

func (c *Consumer) ic() {
	c.emu.State.Patch(types.Patch{
		CursorAt: utils.UintPtr(c.emu.Buf.Addr()),
	})
}

func (c *Consumer) mf(out *Outbound) {
	count := out.MustNext()
	raw := out.MustNextSlice(int(count) * 2)
	cell, _ := c.emu.Buf.Get()
	cell.Attrs = types.NewModifiedAttrs(cell.Attrs, raw)
	c.emu.Buf.SetAndNext(cell)
}

// TODO 🔥 PT not handled
func (c *Consumer) pt() {
	c.emu.Bus.PubPanic("🔥 PT not implemented")
}

func (c *Consumer) ra(out *Outbound, fldAddr uint, fldAttrs *types.Attrs, inFld bool) {
	raw := out.MustNextSlice(2)
	stop := conv.Bytes2Addr(raw)
	// 👇 validate stop addr
	_, ok := c.emu.Buf.Peek(stop)
	if ok {
		char := out.MustNext()
		if types.Order(char) == types.GE {
			// TODO 🔥 GE not properly handled -- what alt character set??
			// also needs to be present in inbound stream (RB, RM/A)
			fldAttrs.LCID = 0xf1
			char = out.MustNext()
		}
		cell := NewCell(c.emu)
		cell.Attrs = fldAttrs
		cell.Char = char
		if inFld {
			cell.SetFldAddr(fldAddr)
		}
		c.emu.Cells.RA(cell, c.emu.Buf.Addr(), stop)
	} else {
		c.emu.Buf.AddrPanic(stop)
	}
}

func (c *Consumer) sa(out *Outbound, fldAttrs *types.Attrs) *types.Attrs {
	chars := out.MustNextSlice(2)
	return types.NewModifiedAttrs(fldAttrs, chars)
}

func (c *Consumer) sba(out *Outbound) {
	raw := out.MustNextSlice(2)
	addr := conv.Bytes2Addr(raw)
	c.emu.Buf.MustSeek(addr)
}

func (c *Consumer) sf(out *Outbound) (uint, *types.Attrs) {
	raw := out.MustNext()
	fldAttrs := types.NewBasicAttrs(raw)
	fldAddr := c.emu.Buf.Addr()
	c.sfImpl(fldAddr, fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfe(out *Outbound) (uint, *types.Attrs) {
	count := out.MustNext()
	raw := out.MustNextSlice(int(count) * 2)
	fldAttrs := types.NewExtendedAttrs(raw)
	fldAddr := c.emu.Buf.Addr()
	c.sfImpl(fldAddr, fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfImpl(fldAddr uint, fldAttrs *types.Attrs) {
	// 🔥 as per spec, if we start a new field at r1/c1 then
	//    treat like an EW -- if we get here after a real EW,
	//    we'll reset a second time -- the clarity of the
	//    code outweighs any small perf hit
	if c.emu.Buf.Addr() == 0 {
		c.emu.Bus.PubReset()
	}
	// 👇 now we can insert the Sf
	sf := NewCell(c.emu)
	sf.Attrs = fldAttrs
	sf.Char = byte(types.SF)
	sf.SetFldAddr(fldAddr)
	c.emu.Buf.SetAndNext(sf)
}

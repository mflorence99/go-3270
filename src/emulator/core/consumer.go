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
	c.emu.Bus.SubInit(c.init)
	c.emu.Bus.SubOutbound(c.consume)
	return c
}

// TODO 🔥 just in case we need it
func (c *Consumer) init() {}

func (c *Consumer) consume(chars []byte) {
	defer utils.ElapsedTime(time.Now())
	// 👇 process the commands in the stream
	out := NewOutbound(chars, c.emu.Bus)
	char := out.MustNext()
	cmd := types.Command(char)
	c.commands(out, cmd)
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
	addr, ok := c.emu.Flds.EAU()
	if ok {
		c.emu.Buf.WrappingSeek(addr + 1)
		c.emu.State.Patch(types.Patch{
			CursorAt: utils.UintPtr(c.emu.Buf.Addr()),
		})
	}
}

func (c *Consumer) ew(out *Outbound) {
	_, ok := c.wcc(out)
	if ok {
		c.emu.Bus.PubReset()
		c.orders(out)
		c.emu.Bus.PubRender()
	}
}

func (c *Consumer) ewa(out *Outbound) {
	_, ok := c.wcc(out)
	if ok {
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

func (c *Consumer) wcc(out *Outbound) (types.WCC, bool) {
	char, ok := out.Next()
	if ok {
		wcc := types.NewWCC(char)
		// TODO 🔥 not yet handled
		if wcc.Reset {
			println("🔥 WCC Reset not implemented")
		}
		if wcc.ResetMDT {
			c.emu.Flds.SetMDTs(false)
		}
		c.emu.Bus.PubWCC(wcc)
		return wcc, true
	} else {
		return types.WCC{}, false
	}
}

// 🟦 WSF (which may contain multiple commands itself)

func (c *Consumer) wsf(out *Outbound) {
	// TODO 🔥 there are a million SF types
	// but we are interested in READ_PARTITION
	sflds := SFldsFromStream(out)
	for _, sfld := range sflds {

		switch sfld.ID {

		case types.READ_PARTITION:
			c.rp(sfld)

		// TODO 🔥 only READ_PARTITION is implemented
		default:
			c.emu.Bus.PubPanic(fmt.Sprintf("🔥 SFld %s not implemented", sfld))

		}
	}
}

func (c *Consumer) rp(sfld SFld) {
	pid := sfld.Info[0]
	if pid == 0xFF {
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

// 🟦 Orders

func (c *Consumer) orders(out *Outbound) {
	var fldAddr uint
	fldAttrs := types.DEFAULT_ATTRS
	// 👇 look at each byte to see if it is an order
outer:
	for out.HasNext() {
		char := out.MustNext()
		order := types.Order(char)
		// 👇 dispatch on order
		switch order {

		// 🔥 per spec invalid EUA terminates write operation
		case types.EUA:
			ok := c.eua(out)
			if !ok {
				break outer
			}

		case types.GE:
			c.ge(out, fldAddr, fldAttrs)

		case types.IC:
			c.ic()

		case types.MF:
			c.mf(out)

		case types.PT:
			c.pt()

		// 🔥 per spec invalid RA terminates write operation
		case types.RA:
			ok := c.ra(out, fldAddr, fldAttrs)
			if !ok {
				break outer
			}

		case types.SA:
			fldAttrs = c.sa(out, fldAttrs)

		case types.SBA:
			c.sba(out)

		case types.SF:
			fldAddr, fldAttrs = c.sf(out)

		case types.SFE:
			fldAddr, fldAttrs = c.sfe(out)

		// 👇 if it isn't an order, it's data
		default:
			c.char(char, fldAddr, fldAttrs)
		}
	}
}

func (c *Consumer) char(char byte, fldAddr uint, fldAttrs *types.Attrs) {
	cell := &buffer.Cell{
		Attrs:   fldAttrs,
		Char:    char,
		FldAddr: fldAddr,
	}
	c.emu.Buf.SetAndNext(cell)
}

func (c *Consumer) eua(out *Outbound) bool {
	raw := out.MustNextSlice(2)
	stop := conv.Bytes2Addr(raw)
	return c.emu.Cells.EUA(c.emu.Buf.Addr(), stop)
}

func (c *Consumer) ge(out *Outbound, fldAddr uint, fldAttrs *types.Attrs) {
	char := out.MustNext()
	// TODO 🔥 GE not properly handled -- what alt character set??
	// also needs to be present in inbound stream (RB, RM/A)
	fldAttrs.LCID = 0xf1
	c.char(char, fldAddr, fldAttrs)
}

func (c *Consumer) ic() {
	cell, _ := c.emu.Buf.Get()
	c.emu.State.Patch(types.Patch{
		CursorAt:  utils.UintPtr(c.emu.Buf.Addr()),
		Numeric:   utils.BoolPtr(cell.Attrs.Numeric),
		Protected: utils.BoolPtr(cell.Attrs.Protected || cell.IsFldStart()),
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

func (c *Consumer) ra(out *Outbound, fldAddr uint, fldAttrs *types.Attrs) bool {
	raw := out.MustNextSlice(2)
	stop := conv.Bytes2Addr(raw)
	char := out.MustNext()
	if types.Order(char) == types.GE {
		// TODO 🔥 GE not properly handled -- what alt character set??
		// also needs to be present in inbound stream (RB, RM/A)
		fldAttrs.LCID = 0xf1
		char = out.MustNext()
	}
	cell := &buffer.Cell{
		Attrs:   fldAttrs,
		Char:    char,
		FldAddr: fldAddr,
	}
	return c.emu.Cells.RA(cell, c.emu.Buf.Addr(), stop)
}

func (c *Consumer) sa(out *Outbound, fldAttrs *types.Attrs) *types.Attrs {
	c.emu.Buf.SetMode(types.CHARACTER_MODE)
	chars := out.MustNextSlice(2)
	return types.NewModifiedAttrs(fldAttrs, chars)
}

func (c *Consumer) sba(out *Outbound) {
	raw := out.MustNextSlice(2)
	addr := conv.Bytes2Addr(raw)
	if addr >= c.emu.Buf.Len() {
		c.emu.Bus.PubPanic("🔥 Data requires a device with a larger screen")
	}
	c.emu.Buf.WrappingSeek(addr)
}

func (c *Consumer) sf(out *Outbound) (uint, *types.Attrs) {
	c.emu.Buf.SetMode(types.FIELD_MODE)
	raw := out.MustNext()
	fldAttrs := types.NewBasicAttrs(raw)
	fldAddr := c.emu.Buf.Addr()
	c.sfImpl(fldAddr, fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfe(out *Outbound) (uint, *types.Attrs) {
	c.emu.Buf.SetMode(types.EXTENDED_FIELD_MODE)
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
	sf := &buffer.Cell{
		Attrs:    fldAttrs,
		Char:     byte(types.SF),
		FldAddr:  fldAddr,
		FldStart: true,
		FldEnd:   false, // 👈 completed by flds.Build()
	}
	c.emu.Buf.SetAndNext(sf)
}

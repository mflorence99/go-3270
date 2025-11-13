package outbound

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/sfld"
	"go3270/emulator/state"
	"go3270/emulator/stream"
	"go3270/emulator/types"
	"go3270/emulator/utils"
	"time"
)

// 🟧 Consume outbound (3270 <- app) data stream

type Consumer struct {
	buf   *buffer.Buffer
	bus   *pubsub.Bus
	cells *buffer.Cells
	cfg   types.Config
	flds  *buffer.Flds
	st    *state.State
}

// 🟦 Constructor

func NewConsumer(bus *pubsub.Bus, buf *buffer.Buffer, cells *buffer.Cells, flds *buffer.Flds, st *state.State) *Consumer {
	c := new(Consumer)
	c.bus = bus
	c.buf = buf
	c.cells = cells
	c.flds = flds
	c.st = st
	// 👇 subscriptions
	c.bus.SubConfig(c.configure)
	c.bus.SubOutbound(c.consume)
	return c
}

func (c *Consumer) configure(cfg types.Config) {
	c.cfg = cfg
}

func (c *Consumer) consume(chars []byte) {
	defer utils.ElapsedTime(time.Now())
	// 👇 process the commands in the stream
	out := stream.NewOutbound(chars, c.bus)
	char := out.MustNext()
	cmd := types.Command(char)
	c.commands(out, cmd)
}

// 🟦 Commands

func (c *Consumer) commands(out *stream.Outbound, cmd types.Command) {
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
	addr := c.flds.EAU()
	if addr != -1 {
		c.buf.WrappingSeek(addr + 1)
		c.st.Patch(state.Patch{
			CursorAt: utils.IntPtr(c.buf.Addr()),
		})
	}
}

func (c *Consumer) ew(out *stream.Outbound) {
	_, ok := c.wcc(out)
	if ok {
		c.bus.PubReset()
		c.orders(out)
		c.bus.PubRender()
	}
}

func (c *Consumer) ewa(out *stream.Outbound) {
	_, ok := c.wcc(out)
	if ok {
		c.bus.PubReset()
		c.orders(out)
		c.bus.PubRender()
	}
}

func (c *Consumer) rb() {
	c.bus.PubRB(types.INBOUND)
}

func (c *Consumer) rm() {
	c.bus.PubRM(types.INBOUND)
}

func (c *Consumer) rma() {
	c.bus.PubRMA(types.INBOUND)
}

func (c *Consumer) w(out *stream.Outbound) {
	c.wcc(out)
	c.orders(out)
	c.bus.PubRender()
}

func (c *Consumer) wcc(out *stream.Outbound) (types.WCC, bool) {
	char, ok := out.Next()
	if ok {
		wcc := types.NewWCC(char)
		// TODO 🔥 not yet handled
		if wcc.Reset {
			println("🔥 WCC Reset not implemented")
		}
		if wcc.ResetMDT {
			c.flds.ResetMDTs()
		}
		c.bus.PubWCC(wcc)
		return wcc, true
	} else {
		return types.WCC{}, false
	}
}

// 🟦 WSF (which may contain multiple commands itself)

func (c *Consumer) wsf(out *stream.Outbound) {
	// TODO 🔥 there are a million SF types
	// but we are interested in READ_PARTITION
	sflds := sfld.SFldsFromStream(out)
	for _, sfld := range sflds {

		switch sfld.ID {

		case types.READ_PARTITION:
			c.rp(sfld)

		// TODO 🔥 only READ_PARTITION is implemented
		default:
			c.bus.PubPanic(fmt.Sprintf("🔥 SFld %s not implemented", sfld))

		}
	}
}

func (c *Consumer) rp(sfld sfld.SFld) {
	pid := sfld.Info[0]
	if pid == 0xFF {
		cmd := sfld.Info[1]

		switch types.Command(cmd) {

		case types.Q:
			c.bus.PubQ()

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
			c.bus.PubQL(qcodes)

		case types.RB:
			c.bus.PubRB(types.INBOUND)

		case types.RM:
			c.bus.PubRM(types.INBOUND)

		case types.RMA:
			c.bus.PubRMA(types.INBOUND)

		}
	}
}

// 🟦 Orders

func (c *Consumer) orders(out *stream.Outbound) {
	fldAddr := -1
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

	// 👇 when the orders have all been processed, build the fields
	c.flds.Build()
}

func (c *Consumer) char(char byte, fldAddr int, fldAttrs *types.Attrs) {
	cell := &buffer.Cell{
		Attrs:   fldAttrs,
		Char:    char,
		FldAddr: fldAddr,
	}
	c.buf.SetAndNext(cell)
}

func (c *Consumer) eua(out *stream.Outbound) bool {
	raw := out.MustNextSlice(2)
	stop := conv.Bytes2Addr(raw)
	return c.cells.EUA(c.buf.Addr(), stop)
}

func (c *Consumer) ge(out *stream.Outbound, fldAddr int, fldAttrs *types.Attrs) {
	char := out.MustNext()
	// TODO 🔥 GE not properly handled -- what alt character set??
	// also needs to be present in inbound stream (RB, RM/A)
	fldAttrs.LCID = 0xf1
	c.char(char, fldAddr, fldAttrs)
}

func (c *Consumer) ic() {
	cell, _ := c.buf.Get()
	c.st.Patch(state.Patch{
		CursorAt:  utils.IntPtr(c.buf.Addr()),
		Numeric:   utils.BoolPtr(cell.Attrs.Numeric),
		Protected: utils.BoolPtr(cell.Attrs.Protected || cell.FldStart),
	})
}

func (c *Consumer) mf(out *stream.Outbound) {
	count := out.MustNext()
	raw := out.MustNextSlice(int(count) * 2)
	cell, _ := c.buf.Get()
	cell.Attrs = types.NewModifiedAttrs(cell.Attrs, raw)
	c.buf.SetAndNext(cell)
}

// TODO 🔥 PT not handled
func (c *Consumer) pt() {
	c.bus.PubPanic("🔥 PT not implemented")
}

func (c *Consumer) ra(out *stream.Outbound, fldAddr int, fldAttrs *types.Attrs) bool {
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
	return c.cells.RA(cell, c.buf.Addr(), stop)
}

func (c *Consumer) sa(out *stream.Outbound, fldAttrs *types.Attrs) *types.Attrs {
	c.buf.SetMode(types.CHARACTER_MODE)
	chars := out.MustNextSlice(2)
	return types.NewModifiedAttrs(fldAttrs, chars)
}

func (c *Consumer) sba(out *stream.Outbound) {
	raw := out.MustNextSlice(2)
	addr := conv.Bytes2Addr(raw)
	if addr >= c.buf.Len() {
		c.bus.PubPanic("🔥 Data requires a device with a larger screen")
	}
	c.buf.WrappingSeek(addr)
}

func (c *Consumer) sf(out *stream.Outbound) (int, *types.Attrs) {
	c.buf.SetMode(types.FIELD_MODE)
	raw := out.MustNext()
	fldAttrs := types.NewBasicAttrs(raw)
	fldAddr := c.buf.Addr()
	c.sfImpl(fldAddr, fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfe(out *stream.Outbound) (int, *types.Attrs) {
	c.buf.SetMode(types.EXTENDED_FIELD_MODE)
	count := out.MustNext()
	raw := out.MustNextSlice(int(count) * 2)
	fldAttrs := types.NewExtendedAttrs(raw)
	fldAddr := c.buf.Addr()
	c.sfImpl(fldAddr, fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfImpl(fldAddr int, fldAttrs *types.Attrs) {
	// 🔥 as per spec, if we start a new field at r1/c1 then
	//    treat like an EW -- if we get here after a real EW,
	//    we'll reset a second time -- the clarity of the
	//    code outweighs any small perf hit
	if c.buf.Addr() == 0 {
		c.bus.PubReset()
	}
	// 👇 now we can insert the Sf
	sf := &buffer.Cell{
		Attrs:    fldAttrs,
		Char:     byte(types.SF),
		FldAddr:  fldAddr,
		FldStart: true,
		FldEnd:   false, // 👈 completed by flds.Build()
	}
	c.buf.SetAndNext(sf)
}

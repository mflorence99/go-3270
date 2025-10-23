package outbound

import (
	"go3270/emulator/attrs"
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/stream"
	"go3270/emulator/utils"
	"go3270/emulator/wcc"
	"time"
)

type Consumer struct {
	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
	st  *state.State
}

func NewConsumer(bus *pubsub.Bus, buf *buffer.Buffer, st *state.State) *Consumer {
	c := new(Consumer)
	c.bus = bus
	c.buf = buf
	c.st = st
	// 👇 subscriptions
	c.bus.SubConfig(c.configure)
	c.bus.SubOutbound(c.consume)
	return c
}

func (c *Consumer) configure(cfg pubsub.Config) {
	c.cfg = cfg
}

func (c *Consumer) consume(chars []byte) {
	defer utils.ElapsedTime(time.Now())
	// 👇 process the commands in the stream
	out := stream.NewOutbound(chars)
	char, _ := out.Next()
	cmd := consts.Command(char)
	c.commands(out, cmd)
}

// 🟦 Commands

func (c *Consumer) commands(out *stream.Outbound, cmd consts.Command) {
	// 👇 dispatch on command
	switch cmd {

	case consts.EAU:
		c.eau()

	case consts.EW:
		c.ew(out)

	case consts.EWA:
		c.ew(out)

	case consts.RB:
		c.rb()

	case consts.RM:
		c.rm()

	case consts.RMA:
		c.rma()

	case consts.W:
		c.w(out)

	case consts.WSF:
		c.wsf(out)
	}
}

func (c *Consumer) eau() {
	c.bus.PubPanic("🔥 EAU not handled")
	// NOTE: EAU doesn't need a WCC
	c.bus.PubRender()
}

func (c *Consumer) ew(out *stream.Outbound) {
	_, ok := c.wcc(out)
	if ok {
		c.bus.PubReset()
		c.orders(out)
		c.normalize()
		c.bus.PubRender()
	}
}

func (c *Consumer) rb() {
	c.bus.PubRM(consts.INBOUND)
}

func (c *Consumer) rm() {
	c.bus.PubRM(consts.INBOUND)
}

func (c *Consumer) rma() {
	c.bus.PubRM(consts.INBOUND)
}

func (c *Consumer) w(out *stream.Outbound) {
	c.wcc(out)
	c.orders(out)
	c.bus.PubRender()
}

func (c *Consumer) wcc(out *stream.Outbound) (wcc.WCC, bool) {
	char, ok := out.Next()
	if ok {
		wcc := wcc.NewWCC(char)
		// TODO 🔥 not yet handled
		if wcc.Reset {
			println("🔥 WCC Reset not implemented")
		}
		if wcc.ResetMDT {
			println("🔥 WCC ResetMDT not implemented")
		}
		c.bus.PubWCC(wcc)
		return wcc, true
	} else {
		return wcc.WCC{}, false
	}
}

func (c *Consumer) wsf(out *stream.Outbound) {
	// 👇 there are a million SF types, but we are interested in READ_PARTITION
	sflds := consts.FromStream(out)
	for _, sfld := range sflds {
		if sfld.ID == consts.READ_PARTITION {
			pid := sfld.Info[0]
			if pid == 0xFF {
				cmd := sfld.Info[1]

				// TODO 🔥 we observe 0xFF when Q is intended!
				if cmd == 0xFF {
					cmd = byte(consts.Q)
				}

				switch consts.Command(cmd) {

				case consts.Q:
					c.bus.PubQ()

				// TODO 🔥 we THINK we can safely return everything
				case consts.QL:
					c.bus.PubQ()

				case consts.RB:
					c.bus.PubRM(consts.INBOUND)

				case consts.RM:
					c.bus.PubRM(consts.INBOUND)

				case consts.RMA:
					c.bus.PubRM(consts.INBOUND)

				}
			}
		}
	}
}

// 🟦 Orders

func (c *Consumer) orders(out *stream.Outbound) {
	fldAddr := 0
	fldAttrs := &attrs.Attrs{Protected: true}
	for out.HasNext() {
		// 👇 look at each byte to see if it is an order
		char, _ := out.Next()
		order := consts.Order(char)
		// 👇 dispatch on order
		switch order {

		case consts.EUA:
			c.eua(out)

		case consts.GE:
			c.ge(out)

		case consts.IC:
			c.ic()

		case consts.MF:
			c.mf(out)

		case consts.PT:
			c.pt()

		case consts.RA:
			c.ra(out)

		case consts.SA:
			fldAttrs = c.sa(out, fldAttrs)

		case consts.SBA:
			c.sba(out)

		case consts.SF:
			fldAddr, fldAttrs = c.sf(out)

		case consts.SFE:
			fldAddr, fldAttrs = c.sfe(out)

		// 👇 if it isn't an order, it's data
		default:
			if char == 0x00 || char >= 0x40 {
				cell := &attrs.Cell{
					Attrs:    fldAttrs,
					Char:     conv.E2A(char),
					FldAddr:  fldAddr,
					FldStart: false,
				}
				c.buf.SetAndNext(cell)
			}
		}
	}
}

func (c *Consumer) eua(out *stream.Outbound) {
	println("🔥 EUA not handled")
	out.NextSlice(2)
}

func (c *Consumer) ge(out *stream.Outbound) {
	println("🔥 GE not handled")
	out.Next()
}

func (c *Consumer) ic() {
	c.st.Patch(state.Patch{
		CursorAt: utils.IntPtr(c.buf.Addr()),
	})
}

func (c *Consumer) mf(out *stream.Outbound) {
	println("🔥 MF not handled")
	count, _ := out.Next()
	out.NextSlice(int(count) * 2)
}

func (c *Consumer) pt() {
	println("🔥 PT not handled")
}

func (c *Consumer) ra(out *stream.Outbound) {
	raw, _ := out.NextSlice(2)
	stop := conv.AddrFromBytes(raw)
	ebcdic, _ := out.Next()
	ascii := conv.E2A(ebcdic)
	// 👇 foundation of what will be repeated
	cell, addr := c.buf.Get()
	attrs := cell.Attrs
	fldAddr := cell.FldAddr
	// 👇 special case: if stop is current address, just fill the buffer
	if stop == addr {
		c.buf.Erase(ascii)
	} else {
		// 👇 watch for wrap around as we blast through to stop
		for {
			cell.Attrs = attrs
			cell.Char = ascii
			cell.FldAddr = fldAddr
			cell.FldStart = false
			cell, addr = c.buf.GetNext()
			if addr == stop {
				break
			}
			c.buf.Seek(addr)
		}
	}
}

func (c *Consumer) sa(out *stream.Outbound, fldAttrs *attrs.Attrs) *attrs.Attrs {
	bytes, _ := out.NextSlice(2)
	return attrs.NewModified(fldAttrs, bytes)
}

func (c *Consumer) sba(out *stream.Outbound) {
	raw, _ := out.NextSlice(2)
	_, ok := c.buf.Seek(conv.AddrFromBytes(raw))
	if !ok {
		c.bus.PubPanic("🔥 Data requires a device with a larger screen")
	}
}

func (c *Consumer) sf(out *stream.Outbound) (int, *attrs.Attrs) {
	next, _ := out.Next()
	fldAttrs := attrs.NewBasic(next)
	fldAddr := c.buf.StartFld(fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfe(out *stream.Outbound) (int, *attrs.Attrs) {
	count, _ := out.Next()
	next, _ := out.NextSlice(int(count) * 2)
	fldAttrs := attrs.NewExtended(next)
	fldAddr := c.buf.StartFld(fldAttrs)
	return fldAddr, fldAttrs
}

// 🟦 Helpers

func (c *Consumer) normalize() {
	// 👇 make sure all cells in a field have the same attributes
	for _, cells := range c.buf.Flds() {
		fld := cells[0]
		for ix := 1; ix < len(cells); ix++ {
			cell := cells[ix]
			cell.Attrs = fld.Attrs
			cell.FldAddr = fld.FldAddr
		}
	}
}

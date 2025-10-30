package outbound

import (
	"fmt"
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
	buf  *buffer.Buffer
	bus  *pubsub.Bus
	cfg  pubsub.Config
	flds *buffer.Flds
	st   *state.State
}

func NewConsumer(bus *pubsub.Bus, buf *buffer.Buffer, flds *buffer.Flds, st *state.State) *Consumer {
	c := new(Consumer)
	c.bus = bus
	c.buf = buf
	c.flds = flds
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

// TODO 🔥 EAU not handled
func (c *Consumer) eau() {
	c.bus.PubPanic("🔥 EAU not handled")
}

func (c *Consumer) ew(out *stream.Outbound) {
	_, ok := c.wcc(out)
	if ok {
		c.bus.PubReset()
		c.orders(out)
		c.flds.Reset()
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
	c.flds.Reset()
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
			c.flds.ResetMDT()
		}
		c.bus.PubWCC(wcc)
		return wcc, true
	} else {
		return wcc.WCC{}, false
	}
}

func (c *Consumer) wsf(out *stream.Outbound) {
	// 👇 there are a million SF types, but we are interested in READ_PARTITION
	sflds := consts.SFldsFromStream(out)
	for _, sfld := range sflds {
		if sfld.ID == consts.READ_PARTITION {
			pid := sfld.Info[0]
			if pid == 0xFF {
				cmd := sfld.Info[1]

				switch consts.Command(cmd) {

				case consts.Q:
					c.bus.PubQ()

				case consts.QL:
					all := (sfld.Info[2] & 0b10000000) == 0x80
					var qcodes []consts.QCode
					if all {
						qcodes = []consts.QCode{
							consts.USABLE_AREA,
							consts.ALPHANUMERIC_PARTITIONS,
							consts.CHARACTER_SETS,
							consts.COLOR_SUPPORT,
							consts.HIGHLIGHTING,
							consts.REPLY_MODES,
							consts.FIELD_VALIDATION,
							consts.FIELD_OUTLINING,
							consts.DDM,
							consts.RPQ_NAMES,
							consts.IMPLICIT_PARTITION,
						}
					} else {
						qcodes = make([]consts.QCode, 0)
						for ix := 3; ix < len(sfld.Info); ix++ {
							qcodes = append(qcodes, consts.QCode(sfld.Info[ix]))
						}
					}
					c.bus.PubQL(qcodes)

				case consts.RB:
					c.bus.PubRB(consts.INBOUND)

				case consts.RM:
					c.bus.PubRM(consts.INBOUND)

				case consts.RMA:
					c.bus.PubRMA(consts.INBOUND)

				}
			}
		}
	}
}

// 🟦 Orders

func (c *Consumer) orders(out *stream.Outbound) {
	fldAddr := 0
	fldAttrs := &attrs.Attrs{Protected: true}
outer:
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

		// 🔥 per spoc invalid RA terminates write operation
		case consts.RA:
			ok := c.ra(out, fldAddr, fldAttrs)
			if !ok {
				break outer
			}

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
				cell := &buffer.Cell{
					Attrs:   fldAttrs,
					Char:    conv.E2A(char),
					FldAddr: fldAddr,
				}
				c.buf.SetAndNext(cell)
			}
		}
	}
}

// TODO 🔥 EUA not handled
func (c *Consumer) eua(out *stream.Outbound) {
	c.bus.PubPanic("🔥 EUA not handled")
	out.NextSlice(2)
}

// TODO 🔥 GE not handled
func (c *Consumer) ge(out *stream.Outbound) {
	c.bus.PubPanic("🔥 GE not handled")
	out.Next()
}

func (c *Consumer) ic() {
	c.st.Patch(state.Patch{
		CursorAt: utils.IntPtr(c.buf.Addr()),
	})
}

// TODO 🔥 MF not handled
func (c *Consumer) mf(out *stream.Outbound) {
	c.bus.PubPanic("🔥 MF not handled")
	count, _ := out.Next()
	out.NextSlice(int(count) * 2)
}

// TODO 🔥 PT not handled
func (c *Consumer) pt() {
	c.bus.PubPanic("🔥 PT not handled")
}

func (c *Consumer) ra(out *stream.Outbound, fldAddr int, fldAttrs *attrs.Attrs) bool {
	raw, _ := out.NextSlice(2)
	stop := conv.AddrFromBytes(raw)
	if stop < c.buf.Len() {
		ebcdic, _ := out.Next()
		ascii := conv.E2A(ebcdic)
		cell, addr := c.buf.Get()
		for {
			if cell == nil {
				cell = buffer.NewCell()
				c.buf.Replace(cell, addr)
			}
			cell.Attrs = fldAttrs
			cell.Char = ascii
			cell.FldAddr = fldAddr
			// 👇 watch for wrap around as we blast through to stop
			cell, addr = c.buf.GetNext()
			if addr == stop {
				c.buf.Seek(stop)
				break
			}
			c.buf.Seek(addr)
		}
		return true
	} else {
		println(fmt.Sprintf("🔥 Inavlid stop address %d in RA order terminates write", stop))
		return false
	}
}

func (c *Consumer) sa(out *stream.Outbound, fldAttrs *attrs.Attrs) *attrs.Attrs {
	c.buf.SetMode(consts.CHARACTER_MODE)
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
	c.buf.SetMode(consts.FIELD_MODE)
	next, _ := out.Next()
	fldAttrs := attrs.NewBasic(next)
	fldAddr := c.buf.StartFld(fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfe(out *stream.Outbound) (int, *attrs.Attrs) {
	c.buf.SetMode(consts.EXTENDED_FIELD_MODE)
	count, _ := out.Next()
	next, _ := out.NextSlice(int(count) * 2)
	fldAttrs := attrs.NewExtended(next)
	fldAddr := c.buf.StartFld(fldAttrs)
	return fldAddr, fldAttrs
}

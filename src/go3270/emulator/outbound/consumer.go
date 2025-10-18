package outbound

import (
	"bytes"
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
	// 👇 dump the stream for debugging
	dmp := pubsub.Dump{
		Bytes:  chars,
		Color:  "yellow",
		EBCDIC: true,
		Title:  "Outbound",
	}
	c.bus.PubDump(dmp)
	// 👇 data can be split into multiple frames
	slices := bytes.Split(chars, consts.LT)
	streams := make([]*stream.Outbound, 0)
	for ix := range slices {
		if len(slices[ix]) > 0 {
			stream := stream.NewOutbound(&slices[ix])
			streams = append(streams, stream)
		}
	}
	// 👇 extract amd process command from each frame
	for _, out := range streams {
		char, ok := out.Next()
		if !ok {
			c.bus.PubPanic("Unable to extract write command")
		}
		cmd := consts.Command(char)
		println(fmt.Sprintf("🐞 Command.%s", cmd))
		c.commands(out, cmd)
	}
	// 👇 render the buffer
	c.bus.PubRender(c.buf.Deltas())
	// 👇 dump the buffer for debugging
	dmp = pubsub.Dump{
		Bytes:  c.buf.Chars(),
		Color:  "plum",
		EBCDIC: false,
		Title:  "Rendered Buffer",
	}
	c.bus.PubDump(dmp)
}

// 🟦 Commands

func (c *Consumer) commands(out *stream.Outbound, cmd consts.Command) {
	// 👇 dispatch on command
	switch cmd {

	case consts.EAU:
		c.eau()

	case consts.EW:
		c.bus.PubReset()
		c.wcc(out)
		c.orders(out)

	case consts.EWA:
		c.bus.PubReset()
		c.wcc(out)
		c.orders(out)

	case consts.RB:
		c.rb()

	case consts.RM:
		c.rm()

	case consts.RMA:
		c.rma()

	case consts.W:
		c.wcc(out)
		c.orders(out)

	case consts.WSF:
		c.wsf()
	}
}

func (c *Consumer) eau() {
	c.bus.PubPanic("🔥 EAU not handled")
}

func (c *Consumer) rb() {
	c.bus.PubPanic("🔥 RB not handled")
}

func (c *Consumer) rm() {
	c.bus.PubInboundRM(consts.INBOUND)
}

func (c *Consumer) rma() {
	c.bus.PubPanic("🔥 RMA not handled")
}

func (c *Consumer) wcc(out *stream.Outbound) {
	char, ok := out.Next()
	if !ok {
		c.bus.PubPanic("Unable to extract WCC")
	}
	wcc := wcc.NewWCC(char)
	println(fmt.Sprintf("🐞 %s", wcc))
	// 👇 honor WCC instructions
	c.st.Patch(state.Patch{
		Alarm:  utils.BoolPtr(wcc.Alarm),
		Locked: utils.BoolPtr(!wcc.Unlock),
	})
	// 🔥 not yet handled
	if wcc.Reset {
		println("🔥 wcc.Reset not handled")
	}
	if wcc.ResetMDT {
		println("🔥 wcc.ResetMDT not handled")
	}
}

func (c *Consumer) wsf() {
	c.bus.PubPanic("🔥 WSF not handled")
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
			c.sa(out)

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

func (c *Consumer) sa(out *stream.Outbound) {
	println("🔥 SA not handled")
	out.NextSlice(2)
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
	if !fldAttrs.Protected {
		println(fmt.Sprintf("🐞 SF at %d %s", c.buf.Addr(), fldAttrs))
	}
	fldAddr := c.buf.StartFld(fldAttrs)
	return fldAddr, fldAttrs
}

func (c *Consumer) sfe(out *stream.Outbound) (int, *attrs.Attrs) {
	count, _ := out.Next()
	next, _ := out.NextSlice(int(count) * 2)
	fldAttrs := attrs.NewExtended(next)
	if !fldAttrs.Protected {
		println(fmt.Sprintf("🐞 SFE at %d %s", c.buf.Addr(), fldAttrs))
	}
	fldAddr := c.buf.StartFld(fldAttrs)
	return fldAddr, fldAttrs
}

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
	"go3270/emulator/stream/outbound"
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
	// 🔥 configure first
	c.bus.SubConfig(c.configure)
	c.bus.SubOutbound(c.consume)
	return c
}

func (c *Consumer) configure(cfg pubsub.Config) {
	c.cfg = cfg
}

func (c *Consumer) consume(chars []byte) {
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
	streams := make([]*outbound.Outbound, 0)
	for ix := range slices {
		if len(slices[ix]) > 0 {
			stream := outbound.NewStream(&slices[ix])
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
	c.bus.PubRender()
	// 👇 dump the buffer for debugging
	dmp = pubsub.Dump{
		Bytes:  c.buf.Chars(),
		Color:  "plum",
		EBCDIC: false,
		Title:  "Rendered Buffer",
	}
	c.bus.PubDump(dmp)
}

func (c *Consumer) commands(out *outbound.Outbound, cmd consts.Command) {
	defer utils.ElapsedTime(time.Now())
	// 👇 dispatch on command
	switch cmd {

	case consts.EAU:
		c.bus.PubPanic("🔥 EAU not handled")

	case consts.EW:
		c.bus.PubReset()
		c.wcc(out)
		c.orders(out)

	case consts.EWA:
		c.bus.PubReset()
		c.wcc(out)
		c.orders(out)

	case consts.RB:
		c.bus.PubPanic("🔥 RB not handled")

	case consts.RM:
		c.bus.PubPanic("🔥 RM not handled")

	case consts.RMA:
		c.bus.PubPanic("🔥 RMA not handled")

	case consts.W:
		c.wcc(out)
		c.orders(out)

	case consts.WSF:
		c.bus.PubPanic("🔥 WSF not handled")
	}
}

func (c *Consumer) orders(out *outbound.Outbound) {
	fldAddr := 0
	fldAttrs := &attrs.Attrs{Protected: true}
	for out.HasNext() {
		// 👇 look at each byte to see if it is an order
		char, _ := out.Next()
		order := consts.Order(char)
		// 👇 dispatch on order
		switch order {

		case consts.EUA:
			println("🔥 EUA not handled")

		case consts.GE:
			println("🔥 GE not handled")

		case consts.IC:
			c.st.Patch(state.Patch{
				CursorAt: utils.IntPtr(c.buf.Addr()),
			})

		case consts.MF:
			println("🔥 MF not handled")

		case consts.PT:
			println("🔥 PT not handled")

		case consts.RA:
			println("🔥 RA not handled")

		case consts.SA:
			println("🔥 SA not handled")

		case consts.SBA:
			raw, _ := out.NextSlice(2)
			_, ok := c.buf.Seek(conv.AddrFromBytes(raw))
			if !ok {
				c.bus.PubPanic("Data requires a device with a larger screen")
			}

		case consts.SF:
			next, _ := out.Next()
			fldAttrs = attrs.NewBasic(next)
			if !fldAttrs.Protected {
				println(fmt.Sprintf("🐞 SF at %d %s", c.buf.Addr(), fldAttrs))
			}
			fldAddr = c.buf.StartFld(fldAttrs)

		case consts.SFE:
			count, _ := out.Next()
			next, _ := out.NextSlice(int(count) * 2)
			fldAttrs = attrs.NewExtended(next)
			if !fldAttrs.Protected {
				println(fmt.Sprintf("🐞 SFE at %d %s", c.buf.Addr(), fldAttrs))
			}
			fldAddr = c.buf.StartFld(fldAttrs)

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

func (c *Consumer) wcc(out *outbound.Outbound) {
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

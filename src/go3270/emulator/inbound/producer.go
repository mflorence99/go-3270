package inbound

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/stream"
)

type Producer struct {
	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
	st  *state.State
}

func NewProducer(bus *pubsub.Bus, buf *buffer.Buffer, st *state.State) *Producer {
	p := new(Producer)
	p.buf = buf
	p.bus = bus
	p.st = st
	// 👇 subscriptions
	p.bus.SubConfig(p.configure)
	p.bus.SubInboundAttn(p.attn)
	p.bus.SubInboundRM(p.rm)
	return p
}

func (p *Producer) attn(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	in.PutSlice(consts.LT)
	title := fmt.Sprintf("Inbound %s", aid)
	p.produce(in.Bytes(), title)
}

func (p *Producer) configure(cfg pubsub.Config) {
	p.cfg = cfg
}

func (p *Producer) produce(bytes []byte, title string) {
	dmp := pubsub.Dump{
		Bytes:  bytes,
		Color:  "palegreen",
		EBCDIC: true,
		Title:  title,
	}
	p.bus.PubDump(dmp)
	p.bus.PubInbound(bytes)
}

func (p *Producer) rm(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	cursorAt := p.st.Stat.CursorAt
	in.PutSlice(conv.AddrToBytes(cursorAt))

	// 🔥 TEMPORARY!!!!

	in.Put(byte(consts.SBA))
	in.PutSlice(conv.AddrToBytes(cursorAt))
	in.Put(conv.A2E('1'))

	in.PutSlice(consts.LT)
	title := fmt.Sprintf("Inbound %s RM", aid)
	p.produce(in.Bytes(), title)
}

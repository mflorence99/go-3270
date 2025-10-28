package inbound

import (
	"go3270/emulator/buffer"
	"go3270/emulator/consts"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/sfld/qr"
	"go3270/emulator/state"
	"go3270/emulator/stream"
)

type Producer struct {
	buf  *buffer.Buffer
	bus  *pubsub.Bus
	cfg  pubsub.Config
	flds *buffer.Flds
	st   *state.State
}

func NewProducer(bus *pubsub.Bus, buf *buffer.Buffer, flds *buffer.Flds, st *state.State) *Producer {
	p := new(Producer)
	p.buf = buf
	p.bus = bus
	p.flds = flds
	p.st = st
	// 👇 subscriptions
	p.bus.SubConfig(p.configure)
	p.bus.SubAttn(p.attn)
	p.bus.SubQ(p.q)
	p.bus.SubRB(p.rb)
	p.bus.SubRM(p.rm)
	p.bus.SubRMA(p.rma)
	return p
}

func (p *Producer) attn(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes())
}

func (p *Producer) configure(cfg pubsub.Config) {
	p.cfg = cfg
}

func (p *Producer) q() {
	in := stream.NewInbound()
	in.Put(byte(consts.INBOUND))
	// 👇 SUMMARY
	qr.NewSummary([]consts.QCode{
		consts.SUMMARY,
		consts.USABLE_AREA,
		consts.ALPHANUMERIC_PARTITIONS,
		consts.CHARACTER_SETS,
		consts.COLOR_SUPPORT,
		consts.HIGHLIGHTING,
		consts.REPLY_MODES,
		consts.RPQ_NAMES,
		consts.DDM,
		// TODO 🔥 this breaks the TERMTEST
		// consts.IMPLICIT_PARTITION,
		consts.FIELD_VALIDATION,
		consts.FIELD_OUTLINING,
	}).Put(in)
	// 👇 then the rest
	qr.NewUsableArea(p.cfg.Cols, p.cfg.Rows).Put(in)
	qr.NewAlphanumericPartitions(p.cfg.Cols, p.cfg.Rows).Put(in)
	qr.NewCharacterSets().Put(in)
	qr.NewColorSupport(p.cfg.Monochrome).Put(in)
	qr.NewHighlighting().Put(in)
	qr.NewReplyModes().Put(in)
	qr.NewRPQNames().Put(in)
	qr.NewDDM().Put(in)
	// TODO 🔥 this breaks the TERMTEST
	// qr.NewImplicitPartition(p.cfg.Cols, p.cfg.Rows).Put(in)
	qr.NewFieldValidation().Put(in)
	qr.NewFieldOutlining().Put(in)
	// 👇 frame boundary LT is last
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes())
}

// TODO 🔥 RB not handled
func (p *Producer) rb(aid consts.AID) {
	p.bus.PubPanic("🔥 RB not handled")
}

func (p *Producer) rm(aid consts.AID) {
	in := stream.NewInbound()
	// 👇 AID + cursor first
	in.Put(byte(aid))
	cursorAt := p.st.Stat.CursorAt
	in.PutSlice(conv.AddrToBytes(cursorAt))
	in.PutSlice(p.flds.ReadMDT())
	// 👇 frame boundary LT is last
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes())
}

// TODO 🔥 RMA not handled
func (p *Producer) rma(aid consts.AID) {
	p.bus.PubPanic("🔥 RMA not handled")
}

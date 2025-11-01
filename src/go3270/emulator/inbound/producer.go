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

// 🟧 Produce inbound (3270 -> app) data stream

type Producer struct {
	buf  *buffer.Buffer
	bus  *pubsub.Bus
	cfg  pubsub.Config
	flds *buffer.Flds
	st   *state.State
}

// 🟦 Constructor

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
	p.bus.SubQL(p.ql)
	p.bus.SubRB(p.rb)
	p.bus.SubRM(p.rm)
	p.bus.SubRMA(p.rma)
	return p
}

func (p *Producer) configure(cfg pubsub.Config) {
	p.cfg = cfg
}

// 🟦 Functions to produce requested stream type

func (p *Producer) attn(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes(), false)
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
		consts.FIELD_VALIDATION,
		consts.FIELD_OUTLINING,
		consts.DDM,
		consts.RPQ_NAMES,
		consts.IMPLICIT_PARTITION,
	}).Put(in)
	// 👇 then the rest
	qr.NewUsableArea(p.cfg.Cols, p.cfg.Rows).Put(in)
	qr.NewAlphanumericPartitions(p.cfg.Cols, p.cfg.Rows).Put(in)
	qr.NewCharacterSets().Put(in)
	qr.NewColorSupport(p.cfg.Monochrome).Put(in)
	qr.NewHighlighting().Put(in)
	qr.NewReplyModes().Put(in)
	qr.NewFieldValidation().Put(in)
	qr.NewFieldOutlining().Put(in)
	qr.NewDDM().Put(in)
	qr.NewRPQNames().Put(in)
	qr.NewImplicitPartition(p.cfg.Cols, p.cfg.Rows).Put(in)
	// 👇 frame boundary LT is last
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes(), true)
}

func (p *Producer) ql(qcodes []consts.QCode) {
	in := stream.NewInbound()
	in.Put(byte(consts.INBOUND))
	for _, qcode := range qcodes {
		switch qcode {
		case consts.USABLE_AREA:
			qr.NewUsableArea(p.cfg.Cols, p.cfg.Rows).Put(in)
		case consts.ALPHANUMERIC_PARTITIONS:
			qr.NewAlphanumericPartitions(p.cfg.Cols, p.cfg.Rows).Put(in)
		case consts.CHARACTER_SETS:
			qr.NewCharacterSets().Put(in)
		case consts.COLOR_SUPPORT:
			qr.NewColorSupport(p.cfg.Monochrome).Put(in)
		case consts.HIGHLIGHTING:
			qr.NewHighlighting().Put(in)
		case consts.REPLY_MODES:
			qr.NewReplyModes().Put(in)
		case consts.FIELD_VALIDATION:
			qr.NewFieldValidation().Put(in)
		case consts.FIELD_OUTLINING:
			qr.NewFieldOutlining().Put(in)
		case consts.DDM:
			qr.NewDDM().Put(in)
		case consts.RPQ_NAMES:
			qr.NewRPQNames().Put(in)
		case consts.IMPLICIT_PARTITION:
			qr.NewImplicitPartition(p.cfg.Cols, p.cfg.Rows).Put(in)
		}
	}
	// 👇 frame boundary LT is last
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes(), true)
}

func (p *Producer) rb(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	cursorAt := p.st.Status.CursorAt
	in.PutSlice(conv.AddrToBytes(cursorAt))
	in.PutSlice(p.flds.ReadBuffer())
	// 👇 frame boundary LT is last
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes(), false)
}

func (p *Producer) rm(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	if !aid.ShortRead() {
		cursorAt := p.st.Status.CursorAt
		in.PutSlice(conv.AddrToBytes(cursorAt))
		in.PutSlice(p.flds.ReadMDTs())
		// 👇 frame boundary LT is last
		in.PutSlice(consts.LT)
		p.bus.PubInbound(in.Bytes(), false)
	}
}

func (p *Producer) rma(aid consts.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	cursorAt := p.st.Status.CursorAt
	in.PutSlice(conv.AddrToBytes(cursorAt))
	in.PutSlice(p.flds.ReadMDTs())
	// 👇 frame boundary LT is last
	in.PutSlice(consts.LT)
	p.bus.PubInbound(in.Bytes(), false)
}

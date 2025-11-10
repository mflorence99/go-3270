package inbound

import (
	"go3270/emulator/buffer"
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/sfld/qr"
	"go3270/emulator/state"
	"go3270/emulator/stream"
	"go3270/emulator/types"
)

// 🟧 Produce inbound (3270 -> app) data stream

type Producer struct {
	buf  *buffer.Buffer
	bus  *pubsub.Bus
	cfg  types.Config
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

func (p *Producer) configure(cfg types.Config) {
	p.cfg = cfg
}

// 🟦 Functions to produce requested stream type

func (p *Producer) attn(aid types.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	in.PutSlice(types.LT)
	p.bus.PubInbound(in.Bytes(), false)
}

func (p *Producer) q() {
	in := stream.NewInbound()
	in.Put(byte(types.INBOUND))
	// 👇 SUMMARY
	qr.NewSummary([]types.QCode{
		types.SUMMARY,
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
	in.PutSlice(types.LT)
	p.bus.PubInbound(in.Bytes(), true)
}

func (p *Producer) ql(qcodes []types.QCode) {
	in := stream.NewInbound()
	in.Put(byte(types.INBOUND))
	for _, qcode := range qcodes {
		switch qcode {
		case types.USABLE_AREA:
			qr.NewUsableArea(p.cfg.Cols, p.cfg.Rows).Put(in)
		case types.ALPHANUMERIC_PARTITIONS:
			qr.NewAlphanumericPartitions(p.cfg.Cols, p.cfg.Rows).Put(in)
		case types.CHARACTER_SETS:
			qr.NewCharacterSets().Put(in)
		case types.COLOR_SUPPORT:
			qr.NewColorSupport(p.cfg.Monochrome).Put(in)
		case types.HIGHLIGHTING:
			qr.NewHighlighting().Put(in)
		case types.REPLY_MODES:
			qr.NewReplyModes().Put(in)
		case types.FIELD_VALIDATION:
			qr.NewFieldValidation().Put(in)
		case types.FIELD_OUTLINING:
			qr.NewFieldOutlining().Put(in)
		case types.DDM:
			qr.NewDDM().Put(in)
		case types.RPQ_NAMES:
			qr.NewRPQNames().Put(in)
		case types.IMPLICIT_PARTITION:
			qr.NewImplicitPartition(p.cfg.Cols, p.cfg.Rows).Put(in)
		}
	}
	// 👇 frame boundary LT is last
	in.PutSlice(types.LT)
	p.bus.PubInbound(in.Bytes(), true)
}

func (p *Producer) rb(aid types.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	cursorAt := p.st.Status.CursorAt
	in.PutSlice(conv.Addr2Bytes(cursorAt))
	in.PutSlice(p.flds.ReadBuffer())
	// 👇 frame boundary LT is last
	in.PutSlice(types.LT)
	p.bus.PubInbound(in.Bytes(), false)
}

func (p *Producer) rm(aid types.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	if !aid.ShortRead() {
		cursorAt := p.st.Status.CursorAt
		in.PutSlice(conv.Addr2Bytes(cursorAt))
		in.PutSlice(p.flds.ReadMDTs())
		// 👇 frame boundary LT is last
		in.PutSlice(types.LT)
		p.bus.PubInbound(in.Bytes(), false)
	}
}

func (p *Producer) rma(aid types.AID) {
	in := stream.NewInbound()
	in.Put(byte(aid))
	cursorAt := p.st.Status.CursorAt
	in.PutSlice(conv.Addr2Bytes(cursorAt))
	in.PutSlice(p.flds.ReadMDTs())
	// 👇 frame boundary LT is last
	in.PutSlice(types.LT)
	p.bus.PubInbound(in.Bytes(), false)
}

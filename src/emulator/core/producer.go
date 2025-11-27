package core

import (
	"emulator/conv"
	"emulator/core/qr"
	"emulator/types"
)

// 🟧 Produce inbound (3270 -> app) data stream

// 👁️ All page references to:
// https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf

type Producer struct {
	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewProducer(emu *Emulator) *Producer {
	p := new(Producer)
	p.emu = emu
	// 👇 subscriptions
	p.emu.Bus.SubAttn(p.attn)
	p.emu.Bus.SubInitialize(p.initialize)
	p.emu.Bus.SubQ(p.q)
	p.emu.Bus.SubQL(p.ql)
	p.emu.Bus.SubRB(p.rb)
	p.emu.Bus.SubRM(p.rm)
	// 🔥 same as RM
	p.emu.Bus.SubRMA(p.rm)
	return p
}

// TODO 🔥 just in case we need it
func (p *Producer) initialize() {}

// 🟦 Functions to produce requested stream type

// 👁️ Short Read Operation p 3-14
func (p *Producer) attn(aid types.AID) {
	in := NewInbound()
	in.Put(byte(aid))
	in.PutSlice(types.LT)
	p.emu.Bus.PubInbound(in.Bytes(), PubInboundHints{Short: true})
}

// 👁️ Query p 6-19
func (p *Producer) q() {
	in := NewInbound()
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
	qr.NewUsableArea(p.emu.Cfg.Cols, p.emu.Cfg.Rows).Put(in)
	qr.NewAlphanumericPartitions(p.emu.Cfg.Cols, p.emu.Cfg.Rows).Put(in)
	qr.NewCharacterSets().Put(in)
	qr.NewColorSupport(p.emu.Cfg.Monochrome).Put(in)
	qr.NewHighlighting().Put(in)
	qr.NewReplyModes().Put(in)
	qr.NewFieldValidation().Put(in)
	qr.NewFieldOutlining().Put(in)
	qr.NewDDM().Put(in)
	qr.NewRPQNames().Put(in)
	qr.NewImplicitPartition(p.emu.Cfg.Cols, p.emu.Cfg.Rows).Put(in)
	// 👇 frame boundary LT is last
	in.PutSlice(types.LT)
	p.emu.Bus.PubInbound(in.Bytes(), PubInboundHints{WSF: true})
}

// 👁️ Query List p 6-19
func (p *Producer) ql(qcodes []types.QCode) {
	in := NewInbound()
	in.Put(byte(types.INBOUND))
	for _, qcode := range qcodes {
		switch qcode {
		case types.USABLE_AREA:
			qr.NewUsableArea(p.emu.Cfg.Cols, p.emu.Cfg.Rows).Put(in)
		case types.ALPHANUMERIC_PARTITIONS:
			qr.NewAlphanumericPartitions(p.emu.Cfg.Cols, p.emu.Cfg.Rows).Put(in)
		case types.CHARACTER_SETS:
			qr.NewCharacterSets().Put(in)
		case types.COLOR_SUPPORT:
			qr.NewColorSupport(p.emu.Cfg.Monochrome).Put(in)
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
			qr.NewImplicitPartition(p.emu.Cfg.Cols, p.emu.Cfg.Rows).Put(in)
		}
	}
	// 👇 frame boundary LT is last
	in.PutSlice(types.LT)
	p.emu.Bus.PubInbound(in.Bytes(), PubInboundHints{WSF: true})
}

// 👁️ Read Buffer command pp 3-12 to 3-13
func (p *Producer) rb(aid types.AID) {
	in := NewInbound()
	in.Put(byte(aid))
	cursorAt := p.emu.State.Status.CursorAt
	in.PutSlice(conv.Addr2Bytes(cursorAt))
	in.PutSlice(p.emu.Cells.RB())
	// 👇 frame boundary LT is last
	in.PutSlice(types.LT)
	p.emu.Bus.PubInbound(in.Bytes(), PubInboundHints{RB: true})
}

// 👁️ Read Modified command pp 3-13 to 3-15
func (p *Producer) rm(aid types.AID) {
	in := NewInbound()
	in.Put(byte(aid))
	if !aid.ShortRead() {
		cursorAt := p.emu.State.Status.CursorAt
		in.PutSlice(conv.Addr2Bytes(cursorAt))
		in.PutSlice(p.emu.Flds.RM())
		// 👇 frame boundary LT is last
		in.PutSlice(types.LT)
		p.emu.Bus.PubInbound(in.Bytes(), PubInboundHints{RM: true})
	}
}

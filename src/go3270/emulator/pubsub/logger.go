package pubsub

import (
	"fmt"
	"go3270/emulator/attrs"
	"go3270/emulator/consts"
	"go3270/emulator/utils"
)

type Logger struct {
	bus *Bus
	cfg Config
}

func NewLogger(bus *Bus) *Logger {
	l := new(Logger)
	l.bus = bus
	// 👇 do this in ctor, so logging precedes actions it logs
	l.bus.SubClose(l.close)
	l.bus.SubConfig(l.configure)
	l.bus.SubDebug(l.debug)
	l.bus.SubInbound(l.produce)
	l.bus.SubOutbound(l.consume)
	l.bus.SubRendered(l.rendered)
	l.bus.SubStatus(l.status)
	l.bus.SubWSF(l.wsf)
	return l
}

func (l *Logger) close() {
	println("🐞 Emulator closed")
}

func (l *Logger) configure(cfg Config) {
	println("🐞 Emulator initialized")
	l.cfg = cfg
}

func (l *Logger) consume(chars []byte) {
	dmp := Dump{
		Bytes:  chars,
		Color:  "yellow",
		EBCDIC: true,
		Title:  "Outbound",
	}
	l.bus.PubDump(dmp)
}

func (l *Logger) debug(topic string, handler interface{}) {
	if topic != "tick" /* 🔥 suppressed ?? && false */ {
		pkg, nm := utils.GetFuncName(handler)
		println(fmt.Sprintf("🐞 topic %s -> func %s() in %s", topic, nm, pkg))
	}
}

func (l *Logger) produce(chars []byte) {
	dmp := Dump{
		Bytes:  chars,
		Color:  "palegreen",
		EBCDIC: true,
		Title:  "Inbound",
	}
	l.bus.PubDump(dmp)
}

func (l *Logger) rendered(chars []byte, flds [][]*attrs.Cell) {
	dmp := Dump{
		Bytes:  chars,
		Color:  "plum",
		EBCDIC: false,
		Title:  "Rendered Buffer",
	}
	l.bus.PubDump(dmp)
	// 👇 log each field in the buffer
	for _, cells := range flds {
		cell := cells[0]
		println(fmt.Sprintf("🐞 SF at %d %s", cell.FldAddr, cell.Attrs))
	}
}

func (l *Logger) status(stat *Status) {
	println(fmt.Sprintf("⚙️ %s", stat))
}

func (l *Logger) wsf(sflds []consts.SFld) {
	println(fmt.Sprintf("🔥 WSF %v", sflds))
}

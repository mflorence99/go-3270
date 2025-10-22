package debug

import (
	"go3270/emulator/buffer"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Logger struct {
	buf *buffer.Buffer
	bus *pubsub.Bus
	cfg pubsub.Config
}

func NewLogger(bus *pubsub.Bus, buf *buffer.Buffer) *Logger {
	l := new(Logger)
	l.buf = buf
	l.bus = bus
	// 👇 do this in ctor, so logging precedes actions it logs
	l.bus.SubClose(l.close)
	l.bus.SubConfig(l.configure)
	l.bus.SubInbound(l.inbound)
	l.bus.SubOutbound(l.outbound)
	l.bus.SubRender(l.render)
	l.bus.SubTrace(l.trace)
	return l
}

func (l *Logger) close() {
	println("🐞 Emulator closed")
}

func (l *Logger) configure(cfg pubsub.Config) {
	l.cfg = cfg
	println("🐞 Emulator initialized")
	LogConfig(l.cfg)
	LogCLUT(l.cfg)
}

func (l *Logger) inbound(chars []byte) {
	LogInbound(chars)
}

func (l *Logger) outbound(chars []byte) {
	LogOutbound(chars)
}

func (l *Logger) render() {
	LogBuffer(l.cfg, l.buf)
}

func (l *Logger) trace(topic string, handler interface{}) {
	LogTrace(topic, handler)
}

// 🟧 Helpers

func Bool(flag bool) string {
	return utils.Ternary(flag, "\u2022", "")
}

func NewTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	return t
}

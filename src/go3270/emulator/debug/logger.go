package debug

import (
	"go3270/emulator/buffer"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"go3270/emulator/wcc"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	l.bus.SubWCC(l.wcc)
	return l
}

func (l *Logger) close() {
	println("🐞 Emulator closed")
}

func (l *Logger) configure(cfg pubsub.Config) {
	l.cfg = cfg
	println("🐞 Emulator initialized")
	l.logConfig(l.cfg)
	l.logCLUT(l.cfg)
}

func (l *Logger) inbound(chars []byte) {
	l.logInbound(chars)
}

func (l *Logger) outbound(chars []byte) {
	l.logOutbound(chars)
}

func (l *Logger) render() {
	l.logBuffer(l.cfg, l.buf)
}

func (l *Logger) trace(topic string, handler interface{}) {
	l.logTrace(topic, handler)
}

func (l *Logger) wcc(wcc wcc.WCC) {
	l.logWCC(wcc)
}

// 🟧 Helpers

func (l *Logger) boolean(flag bool) string {
	return utils.Ternary(flag, "\u2022", "")
}

func (l *Logger) newTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	style := table.StyleBold
	style.Color = table.ColorOptions{
		Border:    text.Colors{text.FgHiBlue, text.Bold},
		Separator: text.Colors{text.FgHiBlue, text.Bold},
		Header:    text.Colors{text.FgHiBlue, text.Bold},
		Row:       text.Colors{text.Reset},
	}
	t.SetStyle(style)
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 1,
			Colors: text.Colors{text.FgHiBlue, text.Bold},
		},
	})
	return t
}

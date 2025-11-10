package logger

import (
	"fmt"
	"go3270/emulator/buffer"
	"go3270/emulator/pubsub"
	"go3270/emulator/state"
	"go3270/emulator/utils"
	"go3270/emulator/wcc"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Subscribe to important messages and log debugging data

// 🔥 Most logging avoids blocking the main thread by using Go routines

type Logger struct {
	buf  *buffer.Buffer
	bus  *pubsub.Bus
	cfg  pubsub.Config
	flds *buffer.Flds
	st   *state.State
}

func NewLogger(bus *pubsub.Bus, buf *buffer.Buffer, flds *buffer.Flds, st *state.State) *Logger {
	l := new(Logger)
	l.buf = buf
	l.bus = bus
	l.flds = flds
	l.st = st
	// 👇 do this in ctor, so logging precedes actions it logs
	l.bus.SubClose(l.close)
	l.bus.SubConfig(l.configure)
	l.bus.SubInbound(l.inbound)
	l.bus.SubOutbound(l.outbound)
	l.bus.SubProbe(l.probe)
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
	go func() {
		l.logConfig()
		l.logCLUT()
	}()
}

func (l *Logger) inbound(chars []byte, wsf bool) {
	go func() {
		// 👇 supplement with an old-fashioned core dump
		l.dump(text.FgHiGreen, "Inbound Core Dump", chars)
		if wsf {
			l.logInboundWSF(chars)
		} else {
			l.logInbound(chars)
		}
	}()
}

func (l *Logger) outbound(chars []byte) {
	go func() {
		// 👇 supplement with an old-fashioned core dump
		l.dump(text.FgYellow, "Outbound Core Dump", chars)
		l.logOutbound(chars)
	}()
}

func (l *Logger) probe(addr int) {
	go func() {
		l.logProbe(addr)
	}()
}

func (l *Logger) render() {
	go func() {
		l.logBuffer()
		l.logFlds()
	}()
}

func (l *Logger) trace(topic string, handler interface{}) {
	// 👇 we need this to be synchronous to make sense
	l.logTrace(topic, handler)
}

func (l *Logger) wcc(wcc wcc.WCC) {
	go func() {
		l.logWCC(wcc)
	}()
}

// 🟧 Helpers

func (l *Logger) boolean(flag bool) string {
	return utils.Ternary(flag, "\u2022", "")
}

func (l *Logger) newTable(color text.Color, title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	style := table.StyleBold
	style.Color = table.ColorOptions{
		Border:    text.Colors{color, text.Bold},
		Separator: text.Colors{color, text.Bold},
		Header:    text.Colors{color, text.Bold},
		Row:       text.Colors{text.Reset},
	}
	t.SetStyle(style)
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 1,
			Colors: text.Colors{color, text.Bold},
		},
	})
	if title != "" {
		t.Style().Title.Align = text.AlignCenter
		t.Style().Title.Colors = text.Colors{color, text.Bold}
		t.SetTitle(title)
	}
	return t
}

func (l *Logger) wrap(w int) text.Transformer {
	return func(val interface{}) string {
		return text.WrapText(fmt.Sprint(val), w)
	}
}

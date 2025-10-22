package pubsub

import (
	"fmt"
	"go3270/emulator/utils"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
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
	l.bus.SubInbound(l.inbound)
	l.bus.SubOutbound(l.outbound)
	return l
}

func (l *Logger) close() {
	println("🐞 Emulator closed")
}

func (l *Logger) configure(cfg Config) {
	println("🐞 Emulator initialized")
	l.cfg = cfg
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "First Name", "Last Name", "Salary"})
	t.AppendRows([]table.Row{
		{1, "Arya", "Stark", 3000},
		{20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
	})
	t.AppendSeparator()
	t.AppendRow([]interface{}{300, "Tyrion", "Lannister", 5000})
	t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.SetStyle(table.StyleLight)
	t.Render()
}

func (l *Logger) debug(topic string, handler interface{}) {
	if topic != "tick" /* 🔥 suppressed ?? */ && false {
		pkg, nm := utils.GetFuncName(handler)
		println(fmt.Sprintf("🐞 topic %s -> func %s() in %s", topic, nm, pkg))
	}
}

func (l *Logger) inbound(chars []byte) {
}

func (l *Logger) outbound(chars []byte) {
}

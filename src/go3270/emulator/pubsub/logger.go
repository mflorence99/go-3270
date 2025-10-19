package pubsub

type Logger struct {
	bus *Bus
	cfg Config
}

func NewLogger(bus *Bus) *Logger {
	l := new(Logger)
	l.bus = bus
	// 👇 subscriptions
	l.bus.SubConfig(l.configure)
	return l
}

func (l *Logger) configure(cfg Config) {
	l.cfg = cfg
	// 👇 ready to log
	l.bus.SubInbound(l.produce)
	l.bus.SubOutbound(l.consume)
	l.bus.SubRendered(l.rendered)
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

func (l *Logger) produce(chars []byte) {
	dmp := Dump{
		Bytes:  chars,
		Color:  "palegreen",
		EBCDIC: true,
		Title:  "Inbound",
	}
	l.bus.PubDump(dmp)
}

func (l *Logger) rendered(chars []byte) {
	dmp := Dump{
		Bytes:  chars,
		Color:  "plum",
		EBCDIC: false,
		Title:  "Rendered Buffer",
	}
	l.bus.PubDump(dmp)
}

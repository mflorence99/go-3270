package logger

import (
	"go3270/emulator/pubsub"
)

// 🟧 Debugger: log inbound (3270 -> app) stream

func (l *Logger) logInbound(chars []byte, hints pubsub.InboundHints) {
	switch {

	case hints.RB:
		l.logInboundRB(chars)

	case hints.RM:
		l.logInboundRM(chars)

	case hints.Short:
		l.logInboundShort(chars)

	case hints.WSF:
		l.logInboundWSF(chars)

	}
}

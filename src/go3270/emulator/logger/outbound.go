package logger

import (
	"go3270/emulator/stream"
	"go3270/emulator/types"

	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: log outbound (3270 <- app) stream

func (l *Logger) logOutbound(chars []byte) {
	// 👇 analyze the commands in the stream
	out := stream.NewOutbound(chars, l.bus)
	char := out.MustNext()
	cmd := types.Command(char)
	// 👇 now we can analyze commands with data
	switch cmd {

	case types.EW:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.EWA:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.W:
		_, ok := out.Next() // 👈 eat the WCC
		if ok {
			l.logOutboundOrders(out, cmd, text.FgYellow)
		}

	case types.WSF:
		l.logOutboundWSF(out, text.FgCyan)
	}
}

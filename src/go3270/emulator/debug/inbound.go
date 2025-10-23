package debug

import (
	"bytes"
	"fmt"
	"go3270/emulator/consts"
	"go3270/emulator/stream"
	"go3270/emulator/utils"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (l *Logger) logInbound(chars []byte) {
	// 👇 convert into a stream for convenience
	slice, _, ok := bytes.Cut(chars, consts.LT)
	out := stream.NewOutbound(utils.Ternary(ok, slice, chars))
	char, _ := out.Next()
	aid := consts.AID(char)
	println(fmt.Sprintf("🐞 %s submitted", aid))
	// 👇 now we can analyze the AID
	switch aid {

	case consts.INBOUND:
		l.logSFlds(out)

	default:
		println(fmt.Sprintf("% 02x", chars[1:]))

	}
}

func (l *Logger) logSFlds(out *stream.Outbound) {
	t := NewTable()
	defer t.Render()
	t.AppendHeader(table.Row{"ID", "QCode", "Info"})
	sflds := consts.FromStream(out)
	for _, sfld := range sflds {
		switch {

		case sfld.ID == consts.QUERY_REPLY:
			qcode := consts.QCode(sfld.Info[0])
			t.AppendRow(table.Row{sfld.ID, qcode, fmt.Sprintf("% 02x", sfld.Info[1:])})

		default:
			t.AppendRow(table.Row{sfld.ID, "", fmt.Sprintf("% 02x", sfld.Info)})

		}
	}
}

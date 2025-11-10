package logger

import (
	"fmt"
	"go3270/emulator/conv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 🟧 Debugger: dump slice of bytes

func (l *Logger) dump(color text.Color, title string, chars []byte) {
	t := l.newTable(color, title)
	defer t.Render()

	// 👇 Header rows

	t.AppendHeader(table.Row{
		"",
		"0        0        0        0        1        1        1        1 ",
		"0                1",
	})

	t.AppendHeader(table.Row{
		"",
		"0 1 2 3  4 5 6 7  8 9 a b  c d e f  0 1 2 3  4 5 6 7  8 9 a b  c d e f",
		"0123456789abcdef 0123456789abcdef",
	})

	// 👇 Dump rows in 32-byte wide slices

	for offset := 0; offset < len(chars); offset += 32 {
		slice := chars[offset:min(offset+32, len(chars))]

		var hex strings.Builder
		for ix := 0; ix < len(slice); ix += 4 {
			chunk := slice[ix:min(ix+4, len(slice))]
			hex.WriteString(fmt.Sprintf("%x ", chunk))
		}

		var ascii strings.Builder
		for ix := 0; ix < len(slice); ix += 16 {
			chunk := slice[ix:min(ix+16, len(slice))]
			// 🔥 we don't know the LCID here, so just dump as ASCII
			ascii.WriteString(fmt.Sprintf("%s ", conv.E2As(string(chunk))))
		}

		t.AppendRow(table.Row{
			fmt.Sprintf("%06x", offset),
			hex.String(),
			ascii.String(),
		})
	}
}

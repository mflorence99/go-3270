package logger

import (
	"fmt"
	"go3270/emulator/types"
)

// 🟦 log Short read

func (l *Logger) logInboundShort(chars []byte) {
	aid := types.AID(chars[0])
	println(fmt.Sprintf("🐞 %s Short Read (3270 -> App)", aid))
}

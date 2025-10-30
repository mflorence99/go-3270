package debug

import (
	"fmt"
	"go3270/emulator/utils"
)

func (l *Logger) logTrace(topic string, handler interface{}) {
	if topic == "rb" /* 🔥 suppressed ?? && false */ {
		pkg, nm := utils.GetFuncName(handler)
		println(fmt.Sprintf("🐞 topic %s -> func %s() in %s", topic, nm, pkg))
	}
}

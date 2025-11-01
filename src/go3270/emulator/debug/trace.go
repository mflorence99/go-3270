package debug

import (
	"fmt"
	"go3270/emulator/utils"
)

// 🟧 Debugger: trace pubsub activity

// TODO 🔥 currently disabled

func (l *Logger) logTrace(topic string, handler interface{}) {
	if topic != "tick" /* 🔥 suppressed ?? */ && false {
		pkg, nm := utils.GetFuncName(handler)
		println(fmt.Sprintf("🐞 topic %s -> func %s() in %s", topic, nm, pkg))
	}
}

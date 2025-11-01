package utils

import (
	"fmt"
	"time"
)

// 🟦 Measure a function's elapsed time (use defer to invode)

func ElapsedTime(start time.Time) {
	elapsed := time.Since(start)
	pkg, nm := GetFuncName(nil)
	println(fmt.Sprintf("⏱️ %s ⏱️ func %s() in %s", elapsed, nm, pkg))
}

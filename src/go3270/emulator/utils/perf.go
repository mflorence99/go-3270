package utils

import (
	"fmt"
	"time"
)

// 👇 defer ElapsedTime(time.Now(), "...") at function start
func ElapsedTime(start time.Time) {
	elapsed := time.Since(start)
	pkg, nm := GetFuncName(nil)
	println(fmt.Sprintf("⏱️ %s -> func %s() in %s", elapsed, pkg, nm))
}

package utils

import (
	"fmt"
	"time"
)

// 🟦 Measure a function's elapsed time (use defer to invode)

func ElapsedTime(start time.Time) {
	elapsed := time.Since(start)
	pkg, nm := GetFuncName(nil)
	fmt.Printf("⏱️ %s ⏱️ func %s() in %s\n", elapsed, nm, pkg)
}

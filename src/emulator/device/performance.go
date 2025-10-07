package device

// 🟧 Measure performance

import (
	"fmt"
	"time"
)

// 👇 defer ElapsedTime(time.Now(), "...") at function start
func ElapsedTime(start time.Time, name string, quiet ...bool) {
	elapsed := time.Since(start)
	if len(quiet) == 0 || !quiet[0] {
		println(fmt.Sprintf("ElapsedTime(%s): %s", name, elapsed))
	}
}

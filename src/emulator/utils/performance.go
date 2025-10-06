package utils

import (
	"fmt"
	"time"
)

// 👇 defer ElapsedTime(time.Now(), "...") at function start
func ElapsedTime(start time.Time, name string, quiet ...bool) {
	elapsed := time.Since(start)
	if len(quiet) == 0 || !quiet[0] {
		fmt.Printf("ElapsedTime(%s): %s\n", name, elapsed)
	}
}

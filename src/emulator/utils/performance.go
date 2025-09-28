package utils

import (
	"fmt"
	"time"
)

// 👇 defer ElapsedTime(time.Now(), "...") at function start
func ElapsedTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

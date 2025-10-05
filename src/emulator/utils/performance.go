package utils

import (
	"fmt"
	"time"
)

type ElapsedTimeOpts struct {
	Quiet bool
}

// 👇 defer ElapsedTime(time.Now(), "...") at function start
func ElapsedTime(start time.Time, name string, opts ElapsedTimeOpts) {
	elapsed := time.Since(start)
	if !opts.Quiet {
		fmt.Printf("ElapsedTime(%s): %s\n", name, elapsed)
	}
}

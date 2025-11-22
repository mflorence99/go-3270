package utils

import (
	"time"
)

func ExampleElapsedTime() {
	defer ElapsedTime(time.Now())
	// ⏱️ 249ns ⏱️ func ExampleElapsedTime() in utils
}

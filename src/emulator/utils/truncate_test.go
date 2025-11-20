package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	str := "Hello, world!"

	trunc := Truncate(str, 50)
	assert.Equal(t, str, trunc, "string shorter than max not truncated")

	short := Truncate(str, 2)
	assert.Equal(t, "He", short, "no ellipses on short truncation")

	long := Truncate(str, 5)
	assert.Equal(t, "He...", long, "ellipses added to longer truncation")
}

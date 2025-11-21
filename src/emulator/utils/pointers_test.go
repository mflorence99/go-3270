package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointers(t *testing.T) {
	b := true
	assert.Equal(t, b, *BoolPtr(b), "bool value equals value from bool ptr")

	s := "Hello, world!"
	assert.Equal(t, s, *StringPtr(s), "string value equals value from string ptr")

	u := uint(255)
	assert.Equal(t, u, *UintPtr(u), "uint value equals value from uint ptr")
}

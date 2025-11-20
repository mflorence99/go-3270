package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenerary(t *testing.T) {
	x := Ternary(true, 1, 2)
	assert.Equal(t, x, 1, "x must be 1")

	y := Ternary(false, 1, 2)
	assert.Equal(t, y, 2, "y must be 2")
}

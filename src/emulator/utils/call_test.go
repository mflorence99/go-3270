package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	c := Call(testCallFunc, 10, 20)
	assert.Equal(t, float64(400), c[0], "call arbitrary function")
}

func testCallFunc(a, b int) float64 {
	return float64(a*b) / 0.5
}

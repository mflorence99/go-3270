package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvert(t *testing.T) {
	p := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	q := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}
	r := Invert(p)
	assert.Equal(t, q, r, "map correectly inverted")
}

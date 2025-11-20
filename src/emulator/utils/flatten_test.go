package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	fat := []any{1, 2, []int{3, 4}, "   str   "}
	flat := Flatten(fat, strings.TrimSpace)
	assert.Equal(t, []any{1, 2, 3, 4, byte(115), byte(116), byte(114)}, flat, "slice flattened")
}

func TestFlatten2Bytes(t *testing.T) {
	fat := []any{1, 2, []int{3, 4}, "   str   "}
	flat := Flatten2Bytes(fat, strings.TrimSpace)
	assert.Equal(t, []byte{byte(1), byte(2), byte(3), byte(4), byte(115), byte(116), byte(114)}, flat, "slice flattened")
}

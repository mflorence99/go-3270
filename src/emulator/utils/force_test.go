package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForceAny2Bytes(t *testing.T) {
	crud := []any{'a', 64, 6.6, "s"}
	bytes := ForceAny2Bytes(crud)
	assert.Equal(t, []byte{0x61, 0x40, 0x06, 0x00}, bytes, "slice reduced to bytes correctly")
}

func TestForceBytes2Any(t *testing.T) {
	bytes := []byte{0x61, 0x40, 0x06}
	crud := ForceBytes2Any(bytes)
	assert.Equal(t, []any{byte(97), byte(64), byte(6)}, crud, "bytes reduced to slice correctly")
}

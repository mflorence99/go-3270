package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLCIDStringer(t *testing.T) {
	assert.Equal(t, "00", LCID(0x00).String(), "LCID 0x00 stringified")
	assert.Equal(t, "f1", LCID(0xf1).String(), "LCID 0xf1 stringified")
}

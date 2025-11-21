package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWCC(t *testing.T) {
	wcc := WCC{
		Alarm:    true,
		Reset:    true,
		ResetMDT: true,
		Unlock:   true,
	}
	assert.Equal(t, wcc.Bits(), byte(0b01000111), "decode WCC to bit settings")
}

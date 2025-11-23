package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {
	emu := MockEmulator().Init()
	assert.Equal(t, uint(0), emu.Buf.Addr(), "initial state")
}

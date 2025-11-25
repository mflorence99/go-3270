package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInbound(t *testing.T) {
	in := NewInbound()
	assert.Empty(t, in.Bytes(), "smoke test passed")
}

func TestInboundPut(t *testing.T) {
	in := NewInbound()
	in.Put('A')
	in.Put('B')
	in.Put('C')
	expected := []byte{'A', 'B', 'C'}
	assert.Equal(t, expected, in.Bytes(), "Put successful")
}

func TestInboundPut16(t *testing.T) {
	in := NewInbound()
	in.Put16(0x1234)
	expected := []byte{0x12, 0x34}
	assert.Equal(t, expected, in.Bytes(), "Put16 successful")
}

func TestInboundSlice(t *testing.T) {
	in := NewInbound()
	expected := []byte{'A', 'B', 'C'}
	in.PutSlice(expected)
	assert.Equal(t, expected, in.Bytes(), "PutSlice successful")
}

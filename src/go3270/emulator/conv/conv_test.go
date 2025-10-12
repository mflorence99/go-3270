package conv_test

import (
	"go3270/emulator/conv"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	assert.True(t, conv.ASCII['0'] == 0xF0)
	assert.True(t, conv.EBCDIC[193-64] == 'A')
}

func Test_A2E(t *testing.T) {
	a := []byte{'H', 'E', 'L', 'L', 'O', ' '}
	e := []byte{200, 197, 211, 211, 214, 64}
	assert.True(t, slices.Equal(e, conv.A2E(a)))
}

func Test_E2A(t *testing.T) {
	a := []byte{'G', 'O', 'O', 'D', 'B', 'Y', 'E', ' '}
	e := []byte{199, 214, 214, 196, 194, 232, 197, 64}
	assert.True(t, slices.Equal(a, conv.E2A(e)))
}

func Test_AddrFromBytes(t *testing.T) {
	addr := 1
	u8s := []byte{0x40, 0xC1}
	assert.True(t, addr == conv.AddrFromBytes(u8s))
}

func Test_AddrToBytes(t *testing.T) {
	addr := 79
	u8s := []byte{0xC1, 0x4F}
	assert.True(t, slices.Equal(u8s, conv.AddrToBytes(addr)))
}

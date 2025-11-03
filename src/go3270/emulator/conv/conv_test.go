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
	assert.True(t, conv.A2E('H') == 200)
}

func Test_E2A(t *testing.T) {
	assert.True(t, conv.E2A(199) == 'G')
}

func Test_AddrFromBytes(t *testing.T) {
	addr := 1
	u8s := []byte{0x40, 0xC1}
	assert.True(t, addr == conv.Bytes2Addr(u8s))
}

func Test_AddrToBytes(t *testing.T) {
	addr := 79
	u8s := []byte{0xC1, 0x4F}
	assert.True(t, slices.Equal(u8s, conv.Addr2Bytes(addr)))
}

package utils_test

import (
	"emulator/utils"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversion_Init(t *testing.T) {
	assert.True(t, utils.ASCII['0'] == 0xF0)
	assert.True(t, utils.EBCDIC[193-64] == 'A')
}

func TestConversion_A2E(t *testing.T) {
	a := []byte{'H', 'E', 'L', 'L', 'O', ' '}
	e := []byte{200, 197, 211, 211, 214, 64}
	assert.True(t, slices.Equal(e, utils.A2E(a)))
}

func TestConversion_E2A(t *testing.T) {
	a := []byte{'G', 'O', 'O', 'D', 'B', 'Y', 'E', ' '}
	e := []byte{199, 214, 214, 196, 194, 232, 197, 64}
	assert.True(t, slices.Equal(a, utils.E2A(e)))
}

func TestConversion_AddrFromBytes(t *testing.T) {
	addr := 1
	bytes := []byte{0x40, 0xC1}
	assert.True(t, addr == utils.AddrFromBytes(bytes))
}

func TestConversion_AddrToBytes(t *testing.T) {
	addr := 79
	bytes := []byte{0xC1, 0x4F}
	assert.True(t, slices.Equal(bytes, utils.AddrToBytes(addr)))
}

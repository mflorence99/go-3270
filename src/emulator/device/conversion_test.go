package device_test

import (
	"emulator/device"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversion_Init(t *testing.T) {
	assert.True(t, device.ASCII['0'] == 0xF0)
	assert.True(t, device.EBCDIC[193-64] == 'A')
}

func TestConversion_A2E(t *testing.T) {
	a := []byte{'H', 'E', 'L', 'L', 'O', ' '}
	e := []byte{200, 197, 211, 211, 214, 64}
	assert.True(t, slices.Equal(e, device.A2E(a)))
}

func TestConversion_E2A(t *testing.T) {
	a := []byte{'G', 'O', 'O', 'D', 'B', 'Y', 'E', ' '}
	e := []byte{199, 214, 214, 196, 194, 232, 197, 64}
	assert.True(t, slices.Equal(a, device.E2A(e)))
}

func TestConversion_AddrFromBytes(t *testing.T) {
	addr := 1
	u8s := []byte{0x40, 0xC1}
	assert.True(t, addr == device.AddrFromBytes(u8s))
}

func TestConversion_AddrToBytes(t *testing.T) {
	addr := 79
	u8s := []byte{0xC1, 0x4F}
	assert.True(t, slices.Equal(u8s, device.AddrToBytes(addr)))
}

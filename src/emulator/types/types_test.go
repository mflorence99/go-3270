package types_test

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes_AID(t *testing.T) {
	assert.True(t, types.AID[0x4C] == "PF24")
	assert.True(t, types.AIDLookup["PF8"] == 0xF8)
}

func TestTypes_CLUT(t *testing.T) {
	assert.True(t, types.CLUT[0xF4] != nil)
	assert.True(t, types.CLUT[0xF8] == nil)
}

func TestTypes_Command(t *testing.T) {
	assert.True(t, types.Command[0xF1] == "W")
	assert.True(t, types.CommandLookup["EW"] == 0xF5)
}

func TestTypes_Op(t *testing.T) {
	assert.True(t, types.Op[0x02] == "Q")
	assert.True(t, types.OpLookup["RM"] == 0xF6)
}

func TestTypes_Order(t *testing.T) {
	assert.True(t, types.Order[0x1D] == "SF")
	assert.True(t, types.OrderLookup["EUA"] == 0x12)
}

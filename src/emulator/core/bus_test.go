package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBus(t *testing.T) {
	bus := NewBus()
	bus.SubInitialize(func() {
		assert.Fail(t, "SubInit() should not be called here")
	})
	expected := []byte{0x01, 0x02, 0x03}
	bus.SubOutbound(func(actual []byte) {
		assert.Equal(t, expected, actual, "SubOutbound() received correctly")
	})
	bus.PubOutbound(expected)
}

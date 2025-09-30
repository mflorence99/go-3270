package device_test

import (
	"emulator/device"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestTest(t *testing.T) {
	_ = device.NewDevice(nil, "", 0, nil, 0, 0, 0, 0, 0, 0, 0)
	assert.True(t, true)
}

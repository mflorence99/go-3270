package device_test

import (
	"emulator/device"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevice_frames(t *testing.T) {
	d := device.NewDevice(nil, nil, "", "", 0, 0, 0, 0, 0, 0, 0, 0)
	frames := d.MakeFramesFromBytes([]uint8{0x00, 0x00, 0xFF, 0xEF, 0x00, 0x00, 0xFF, 0xEF, 0x00})
	assert.True(t, len(frames) == 3)
}

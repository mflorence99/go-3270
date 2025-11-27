package core

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEmulator(t *testing.T) {
	emu := MockEmulator(12, 40)
	ok := false
	emu.Bus.SubRender(func() {
		ok = true
	})
	emu.Initialize()
	stream := MockStream(types.EW, types.WCC{}, MockExampleImg, MockExampleAttrs)
	emu.Bus.PubOutbound(stream)
	assert.True(t, ok, "smoke test for mock render")
}

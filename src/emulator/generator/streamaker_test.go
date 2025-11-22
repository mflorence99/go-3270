package generator

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreammaker(t *testing.T) {
	emu := MockEmulator()
	emu.Bus.SubRender(func() {
		assert.True(t, true, "smoke test for mock render")
	})
	emu.Init()
	stream := MakeStream(types.EW, types.WCC{}, ExampleImg, ExampleAttrs)
	emu.Bus.PubOutbound(stream)
}

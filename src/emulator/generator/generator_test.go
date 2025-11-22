package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockGenerator(t *testing.T) {
	emu := MockEmulator()
	emu.Bus.SubInit(func() {
		assert.True(t, true, "smoke test for mock  init")
	})
	emu.Init()
}

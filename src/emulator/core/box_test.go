package core

import (
	"emulator/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBox(t *testing.T) {
	cfg := types.Config{
		FontHeight:   16,
		FontSize:     12,
		FontWidth:    9,
		PaddedHeight: 1.5,
		PaddedWidth:  1.1,
	}
	box := NewBox(5, 10, &cfg)
	assert.Equal(t, 90.0, box.X)
	assert.Equal(t, 96.0, box.Y)
	assert.Equal(t, 10.0, box.W)
	assert.Equal(t, 24.0, box.H)
}

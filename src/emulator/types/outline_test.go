package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutlineStringer(t *testing.T) {
	outline := Outline(0b00001111)
	assert.Equal(t, "BRTL", outline.String(), "outline stringified")
	assert.Equal(t, "BRTL", OutlineFor(outline), "outline stringified")
}

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandStringer(t *testing.T) {
	assert.Equal(t, "WSF", WSF.String(), "WSF stringified")
	assert.Equal(t, "WSF", CommandFor(WSF), "WSF stringified")
}

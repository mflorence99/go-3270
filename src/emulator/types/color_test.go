package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorStringer(t *testing.T) {
	assert.Equal(t, "YELLOW", YELLOW.String(), "YELLOW stringified")
	assert.Equal(t, "YELLOW", ColorFor(YELLOW), "YELLOW stringified")
}

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModeStringer(t *testing.T) {
	assert.Equal(t, "EXTENDED_FIELD_MODE", EXTENDED_FIELD_MODE.String(), "EXTENDED_FIELD_MODE stringified")
	assert.Equal(t, "EXTENDED_FIELD_MODE", ModeFor(EXTENDED_FIELD_MODE), "EXTENDED_FIELD_MODE stringified")
}

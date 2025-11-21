package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighlightStringer(t *testing.T) {
	assert.Equal(t, "UNDERSCORE", UNDERSCORE.String(), "UNDERSCORE stringified")
	assert.Equal(t, "UNDERSCORE", HighlightFor(UNDERSCORE), "UNDERSCORE stringified")
}

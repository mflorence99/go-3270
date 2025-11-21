package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQCodeStringer(t *testing.T) {
	assert.Equal(t, "RPQ_NAMES", RPQ_NAMES.String(), "RPQ_NAMES stringified")
	assert.Equal(t, "RPQ_NAMES", QCodeFor(RPQ_NAMES), "RPQ_NAMES stringified")
}

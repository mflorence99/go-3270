package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStringer(t *testing.T) {
	assert.Equal(t, "SBA", SBA.String(), "SBA stringified")
	assert.Equal(t, "SBA", OrderFor(SBA), "SBA stringified")
}

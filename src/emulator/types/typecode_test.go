package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypecodeStringer(t *testing.T) {
	assert.Equal(t, "CHARSET", CHARSET.String(), "CHARSET stringified")
	assert.Equal(t, "CHARSET", TypecodeFor(CHARSET), "CHARSET stringified")
}

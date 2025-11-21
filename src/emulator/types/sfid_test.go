package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSFIDStringer(t *testing.T) {
	assert.Equal(t, "QUERY_REPLY", QUERY_REPLY.String(), "QUERY_REPLY stringified")
	assert.Equal(t, "QUERY_REPLY", SFIDFor(QUERY_REPLY), "QUERY_REPLY stringified")
}

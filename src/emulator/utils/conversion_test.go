package utils_test

import (
	"emulator/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	assert.True(t, utils.ASCII['0'] == 0xF0)
	assert.True(t, utils.EBCDIC[193-64] == 'A')
}

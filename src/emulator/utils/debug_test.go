package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuncName(t *testing.T) {
	pkg, nm := GetFuncName(anArbitraryFunction)
	assert.Equal(t, "utils", pkg, "extract package name from function")
	assert.Equal(t, "anArbitraryFunction", nm, "extract function name from function")
}

func anArbitraryFunction() {}

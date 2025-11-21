package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	stack := NewStack[uint](8)

	var ix uint
	var ok bool

	assert.True(t, stack.Empty(), "stack is empty initially")
	assert.Equal(t, 0, stack.Len(), "stack has zero length initially")
	_, ok = stack.Peek()
	assert.True(t, !ok, "stack has no items to Peek initially")
	_, ok = stack.Pop()
	assert.True(t, !ok, "stack has no items to Pop initially")

	stack.Push(uint(10))

	assert.False(t, stack.Empty(), "stack is not empty after Push")
	assert.Equal(t, 1, stack.Len(), "stack has length after Push")
	ix, ok = stack.Peek()
	assert.True(t, ok, "stack has items to Peek after Push")
	assert.Equal(t, uint(10), ix, "stack has items to Peek after Push")
	ix, ok = stack.Pop()
	assert.True(t, ok, "stack has items to Pop after Push")
	assert.Equal(t, uint(10), ix, "stack has items to Pop after Push")

	assert.True(t, stack.Empty(), "stack is empty after Pop")
}

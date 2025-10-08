package stack_test

import (
	"emulator/stack"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Stack(t *testing.T) {
	s := stack.New[int](10)
	assert.True(t, s.Len() == 0)
	assert.True(t, s.Empty())
	s.Push(1)
	s.Push(2)
	v, ok := s.Pop()
	assert.True(t, s.Len() == 1 && v == 2 && ok)
	v, ok = s.Peek()
	assert.True(t, s.Len() == 1 && v == 1 && ok)
	v, ok = s.Pop()
	assert.True(t, s.Len() == 0 && v == 1 && ok)
	v, ok = s.Pop()
	assert.True(t, s.Len() == 0 && v == 0 && !ok)
}

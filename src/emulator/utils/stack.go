package utils

type Stack[T any] struct {
	data []T
}

func NewStack[T any](capacity int) *Stack[T] {
	return &Stack[T]{data: make([]T, 0, capacity)}
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) Peek() T {
	if len(s.data) == 0 {
		panic("peek from empty stack")
	}
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) Pop() T {
	if len(s.data) == 0 {
		panic("pop from empty stack")
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

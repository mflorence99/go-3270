package utils

// 🟧 Simple stack implementation

type Stack[T any] struct {
	data []T
}

// 🟦 Constructor

func NewStack[T any](capacity int) *Stack[T] {
	return &Stack[T]{data: make([]T, 0, capacity)}
}

// 🟦 Public functions

func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		return *new(T), false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		return *new(T), false
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, true
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

package util

import a "github.com/kode4food/sputter/api"

// Stack is your standard Stack interface
type Stack interface {
	Push(v a.Value)
	Peek() (v a.Value, ok bool)
	Pop() (v a.Value, ok bool)
}

// basicStack implements a non-concurrent Stack
type basicStack struct {
	head *entry
}

type entry struct {
	value a.Value
	next  *entry
}

// NewStack creates a new Stack instance
func NewStack() Stack {
	return &basicStack{}
}

// Push a Value onto the Stack
func (s *basicStack) Push(v a.Value) {
	if s.head == nil {
		s.head = &entry{v, nil}
		return
	}
	s.head = &entry{v, s.head}
}

// Peek the head of the Stack
func (s *basicStack) Peek() (a.Value, bool) {
	e := s.head
	if e != nil {
		return e.value, true
	}
	return nil, false
}

// Pop the head of the Stack
func (s *basicStack) Pop() (a.Value, bool) {
	e := s.head
	if e != nil {
		s.head = e.next
		return e.value, true
	}
	return nil, false
}

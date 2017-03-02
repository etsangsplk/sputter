package util

// value allows the Stack to take anything
type value interface {
}

// Stack is your standard Stack interface
type Stack interface {
	Push(v value)
	Peek() (v value, ok bool)
	Pop() (v value, ok bool)
}

// basicStack implements a non-concurrent Stack
type basicStack struct {
	head *entry
}

type entry struct {
	value value
	next  *entry
}

// NewStack creates a new Stack instance
func NewStack() Stack {
	return &basicStack{}
}

// Push a Value onto the Stack
func (s *basicStack) Push(v value) {
	if s.head == nil {
		s.head = &entry{v, nil}
		return
	}
	s.head = &entry{v, s.head}
}

// Peek the head of the Stack
func (s *basicStack) Peek() (value, bool) {
	e := s.head
	if e != nil {
		return e.value, true
	}
	return nil, false
}

// Pop the head of the Stack
func (s *basicStack) Pop() (value, bool) {
	e := s.head
	if e != nil {
		s.head = e.next
		return e.value, true
	}
	return nil, false
}

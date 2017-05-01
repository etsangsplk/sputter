package util

// Any allows the utilities to take anything
type Any interface{}

// Stack is your standard Stack interface
type Stack interface {
	Push(value Any)
	Peek() (Any, bool)
	Pop() (Any, bool)
}

// stack implements a non-concurrent Stack
type stack struct {
	head *entry
}

type entry struct {
	value Any
	next  *entry
}

// NewStack creates a new Stack instance
func NewStack() Stack {
	return &stack{}
}

// Push a Value onto the Stack
func (s *stack) Push(value Any) {
	if s.head == nil {
		s.head = &entry{
			value: value,
			next:  nil,
		}
		return
	}
	s.head = &entry{
		value: value,
		next:  s.head,
	}
}

// Peek the head of the Stack
func (s *stack) Peek() (Any, bool) {
	e := s.head
	if e != nil {
		return e.value, true
	}
	return nil, false
}

// Pop the head of the Stack
func (s *stack) Pop() (Any, bool) {
	e := s.head
	if e != nil {
		s.head = e.next
		return e.value, true
	}
	return nil, false
}

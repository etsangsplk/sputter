package api

// Stack implements a non-concurrent stack
type Stack struct {
	head *entry
}

type entry struct {
	value Value
	next  *entry
}

// Push a Value onto the Stack
func (s *Stack) Push(v Value) {
	if s.head == nil {
		s.head = &entry{v, nil}
		return
	}
	s.head = &entry{v, s.head}
}

// Peek the head of the Stack
func (s *Stack) Peek() (v Value, ok bool) {
	e := s.head
	if e != nil {
		return e.value, true
	}
	return nil, false
}

// Pop the head of the Stack
func (s *Stack) Pop() (v Value, ok bool) {
	e := s.head
	if e != nil {
		s.head = e.next
		return e.value, true
	}
	return nil, false
}

package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

const defaultOperandStackSize = 8

type operandStack struct {
	data []interface{}
	SP   int
}

func newOperandStack(size int) *operandStack {
	return &operandStack{
		data: make([]interface{}, size),
		SP:   size,
	}
}

func (s *operandStack) push(v interface{}) {
	s.SP--
	s.data[s.SP] = v
}

func (s *operandStack) pop() interface{} {
	r := s.data[s.SP]
	s.SP++
	return r
}

func (s *operandStack) peek() interface{} {
	return s.data[s.SP]
}

func Exec(c a.Context, bc []byte) a.Value {
	var PC uint
	var i int
	var v interface{}

	stack := newOperandStack(defaultOperandStackSize)

start:
	switch bc[PC] {
	case NoOp:
	case PushFrame:
		i = stack.pop().(int)
		v = stack
		stack = newOperandStack(i)
		stack.push(v)
	case PopFrame:
		stack = stack.pop().(*operandStack)
	case PushContext:
		stack.push(c)
		c = a.ChildContext(c)
	case PopContext:
		c = stack.pop().(a.Context)
	case Nil:
		stack.push(a.Nil)
	case Return:
		return stack.pop().(a.Value)
	case Halt:
		panic(stack.pop().(a.Value))
	}
	PC++
	goto start
}

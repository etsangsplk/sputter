package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

type (
	// Data represents the data segment of a Module
	Data []a.Value

	// Instruction represents a decoded VM instruction
	Instruction struct {
		OpCode OpCode
		Arg0   uint
	}

	// Module is the basic translation unit for the VM
	Module struct {
		a.BaseFunction
		LocalsSize   uint
		StackSize    uint
		Data         []a.Value
		Instructions []Instruction
	}

	operandStack []a.Value
)

// Apply makes Module applicable
func (m *Module) Apply(c a.Context, args a.Sequence) a.Value {
	// Registers
	var r1 a.Value
	var u1, u2 uint
	var s1 a.Sequence
	var PC uint
	var SP uint

	LOCALS := make([]a.Value, m.LocalsSize)
	LOCALS[0] = args

	STACK := make(operandStack, m.StackSize)
	SP = m.StackSize - 1

	DATA := m.Data
	INST := m.Instructions

	push := func(v a.Value) {
		STACK[SP] = v
		SP--
	}

	pop := func() a.Value {
		SP++
		return STACK[SP]
	}

start:
	switch INST[PC].OpCode {
	case NoOp:
		PC++
		goto start

	case Load:
		push(LOCALS[INST[PC].Arg0])
		PC++
		goto start

	case Store:
		LOCALS[INST[PC].Arg0] = pop()
		PC++
		goto start

	case Clear:
		LOCALS[INST[PC].Arg0] = a.Nil
		PC++
		goto start

	case Dup:
		STACK[SP] = STACK[SP+1]
		SP--
		PC++
		goto start

	case Swap:
		u1 = SP + 1
		u2 = SP + 2
		r1 = STACK[u1]
		STACK[u1] = STACK[u2]
		STACK[u2] = r1
		r1 = nil // gc
		PC++
		goto start

	case Nil:
		push(a.Nil)
		PC++
		goto start

	case EmptyList:
		push(a.EmptyList)
		PC++
		goto start

	case True:
		push(a.True)
		PC++
		goto start

	case False:
		push(a.False)
		PC++
		goto start

	case Zero:
		push(a.Zero)
		PC++
		goto start

	case One:
		push(a.One)
		PC++
		goto start

	case Const:
		push(DATA[INST[PC].Arg0])
		PC++
		goto start

	case Def:
		r1 = pop()
		a.GetContextNamespace(c).Put(a.AssertUnqualified(pop()).Name(), r1)
		r1 = nil // gc
		PC++
		goto start

	case Let:
		r1 = pop()
		c.Put(a.AssertUnqualified(pop()).Name(), r1)
		r1 = nil // gc
		PC++
		goto start

	case Eval:
		push(a.Eval(c, pop()))
		PC++
		goto start

	case Apply:
		s1 = a.AssertSequence(pop())
		push(a.AssertApplicable(pop()).Apply(c, s1))
		s1 = nil // gc
		PC++
		goto start

	case First:
		push(a.AssertSequence(pop()).First())
		PC++
		goto start

	case Rest:
		push(a.AssertSequence(pop()).Rest())
		PC++
		goto start

	case Split:
		r1, s1 = a.AssertSequence(pop()).Split()
		push(s1)
		push(r1)
		r1 = nil // gc
		s1 = nil // gc
		PC++
		goto start

	case Prepend:
		r1 = pop()
		push(a.AssertSequence(pop()).Prepend(r1))
		r1 = nil // gc
		PC++
		goto start

	case CondJump:
		r1 = pop()
		if r1 == a.False || r1 == a.Nil {
			r1 = nil // gc
			PC++
			goto start
		}
		r1 = nil // gc
		PC = INST[PC].Arg0
		goto start

	case Jump:
		PC = INST[PC].Arg0
		goto start

	case Return:
		return pop()

	case ReturnNil:
		return a.Nil

	case Panic:
		panic(pop())
	}

	panic("how did we get here?")
}

// WithMetadata creates a copy of this Module with additional Metadata
func (m *Module) WithMetadata(md a.Object) a.AnnotatedValue {
	return &Module{
		BaseFunction: m.Extend(md),
		LocalsSize:   m.LocalsSize,
		StackSize:    m.StackSize,
		Data:         m.Data,
		Instructions: m.Instructions,
	}
}

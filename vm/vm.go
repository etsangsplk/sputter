package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

type (
	// Data represents the data segment of a Module
	Data []a.Value

	// Instruction represents a decoded VM instruction
	Instruction struct {
		OpCode OpCode
		Args   [3]uint
	}

	Module struct {
		LocalsSize   uint
		StackSize    uint
		Data         []a.Value
		Instructions []Instruction
	}

	operandStack []interface{}
)

func (m *Module) Apply(c a.Context, args a.Sequence) a.Value {
	var emptyName a.Name

	// Registers
	var r1 interface{}
	var n1 a.Name
	var a1 a.Applicable
	var s1 a.Sequence
	var PC uint
	var SP uint

	LOCALS := make([]a.Value, m.LocalsSize)
	LOCALS[0] = args

	STACK := make(operandStack, m.StackSize)
	SP = m.StackSize - 1

	DATA := m.Data
	INST := m.Instructions

	push := func(v interface{}) {
		STACK[SP] = v
		SP--
	}

	pop := func() interface{} {
		SP++
		return STACK[SP]
	}

	peek := func() interface{} {
		return STACK[SP+1]
	}

start:
	switch INST[PC].OpCode {
	case NoOp:
		PC++
		goto start

	case Load:
		push(LOCALS[INST[PC].Args[0]])
		PC++
		goto start

	case Store:
		LOCALS[INST[PC].Args[0]] = pop().(a.Value)
		PC++
		goto start

	case Dup:
		push(peek())
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
		push(DATA[INST[PC].Args[0]])
		PC++
		goto start

	case PushContext:
		push(c)
		c = a.ChildContext(c)
		PC++
		goto start

	case PopContext:
		c = pop().(a.Context)
		PC++
		goto start

	case Def:
		n1 = a.AssertUnqualified(pop().(a.Value)).Name()
		a.GetContextNamespace(c).Put(n1, pop().(a.Value))
		n1 = emptyName // gc
		PC++
		goto start

	case Let:
		n1 = a.AssertUnqualified(pop().(a.Value)).Name()
		c.Put(n1, pop().(a.Value))
		n1 = emptyName // gc
		PC++
		goto start

	case Eval:
		push(a.Eval(c, pop().(a.Value)))
		PC++
		goto start

	case Apply:
		a1 = a.AssertApplicable(pop().(a.Value))
		push(a1.Apply(c, a.AssertSequence(pop().(a.Value))))
		a1 = nil // gc
		PC++
		goto start

	case First:
		push(a.AssertSequence(pop().(a.Value)).First())
		PC++
		goto start

	case Rest:
		push(a.AssertSequence(pop().(a.Value)).Rest())
		PC++
		goto start

	case Prepend:
		s1 = a.AssertSequence(pop().(a.Value))
		push(s1.Prepend(pop().(a.Value)))
		PC++
		goto start

	case Truthy:
		r1 = pop()
		if r1 == a.False || r1 == a.Nil {
			push(0)
			r1 = nil // gc
			PC++
			goto start
		}
		push(1)
		r1 = nil // gc
		PC++
		goto start

	case CondJump:
		PC = INST[PC].Args[pop().(uint)]
		goto start

	case Jump:
		PC = INST[PC].Args[0]
		goto start

	case Return:
		return pop().(a.Value)

	case ReturnNil:
		return a.Nil

	case Panic:
		panic(pop())
	}

	panic("how did we get here?")
}

func (m *Module) Str() a.Str {
	return a.MakeDumpStr(m)
}

package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

type (
	// Data represents the data segment of a Module
	Data a.Values

	// Instruction represents a decoded VM instruction
	Instruction struct {
		OpCode OpCode
		Op1    uint
		Op2    uint
	}

	// Module is the basic translation unit for the VM
	Module struct {
		a.BaseFunction
		LocalsSize   uint
		StackSize    uint
		Data         a.Values
		Instructions []Instruction
	}
)

var negOne = a.Zero.Sub(a.One)

// Apply makes Module applicable
func (m *Module) Apply(c a.Context, args a.Sequence) a.Value {
	// Registers
	var c1 a.Context
	var r1 a.Value
	var u1, u2 uint
	var s1 a.Sequence
	var v1 a.Values
	var e1 a.Evaluable
	var b1 bool
	var PC uint
	var SP uint

	LOCALS := make(a.Values, m.LocalsSize)
	LOCALS[Context] = c
	LOCALS[Args] = args

	STACK := make(a.Values, m.StackSize)
	SP = m.StackSize - 1

	DATA := m.Data
	INST := m.Instructions

	push := func(v a.Value) {
		STACK[SP] = v
		SP--
	}

	pop := func() a.Value {
		// TODO: require a clean way to clear cells to avoid leaking
		SP++
		return STACK[SP]
	}

start:
	switch INST[PC].OpCode {
	case NoOp:
		PC++
		goto start

	case Pop:
		pop()
		PC++
		goto start

	case Load:
		push(LOCALS[INST[PC].Op1])
		PC++
		goto start

	case Store:
		LOCALS[INST[PC].Op1] = pop()
		PC++
		goto start

	case StoreConst:
		LOCALS[INST[PC].Op2] = DATA[INST[PC].Op1]
		PC++
		goto start

	case Clear:
		LOCALS[INST[PC].Op1] = a.Nil
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

	case NegOne:
		push(negOne)
		PC++
		goto start

	case Const:
		push(DATA[INST[PC].Op1])
		PC++
		goto start

	case Def:
		r1 = pop()
		a.GetContextNamespace(c).Put(a.AssertUnqualified(pop()).Name(), r1)
		r1 = nil // gc
		PC++
		goto start

	case Let:
		c1 = LOCALS[INST[PC].Op1].(a.Context)
		r1 = pop()
		c1.Put(a.AssertUnqualified(pop()).Name(), r1)
		c1 = nil // gc
		r1 = nil // gc
		PC++
		goto start

	case Eval:
		r1, _ = a.MacroExpand(c, pop())
		if e1, b1 = r1.(a.Evaluable); b1 {
			push(e1.Eval(c))
			e1 = nil // gc
		} else {
			push(r1)
		}
		r1 = nil // gc
		PC++
		goto start

	case Apply:
		s1 = a.AssertSequence(pop())
		push(a.AssertApplicable(pop()).Apply(c, s1))
		s1 = nil // gc
		PC++
		goto start

	case Call:
		u1 = INST[PC].Op1
		u2 = SP + u1 + 1
		v1 = make(a.Values, u1)
		copy(v1, STACK[SP+1:u2])
		STACK[u2] = STACK[u2].(a.Applicable).Apply(c, v1)
		SP = u2 - 1
		v1 = nil // gc
		PC++
		goto start

	case Vector:
		u1 = INST[PC].Op1
		u2 = SP + u1
		v1 = make(a.Values, u1)
		copy(v1, STACK[SP+1:u2+1])
		STACK[u2] = v1
		SP = u2 - 1
		v1 = nil // gc
		PC++
		goto start

	case IsSeq:
		if s1, b1 = pop().(a.Sequence); b1 && s1.IsSequence() {
			push(a.True)
		} else {
			push(a.False)
		}
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
		if s1, b1 = pop().(a.Sequence); !b1 {
			push(a.False)
			PC++
			goto start
		}
		if r1, s1, b1 = s1.Split(); b1 {
			push(s1)
			push(r1)
			push(a.True)
			r1 = nil // gc
			s1 = nil // gc
		} else {
			push(a.False)
		}
		PC++
		goto start

	case Prepend:
		r1 = pop()
		push(a.AssertSequence(pop()).Prepend(r1))
		r1 = nil // gc
		PC++
		goto start

	case Add:
		push(pop().(a.Number).Add(pop().(a.Number)))
		PC++
		goto start

	case Sub:
		push(pop().(a.Number).Sub(pop().(a.Number)))
		PC++
		goto start

	case Mul:
		push(pop().(a.Number).Mul(pop().(a.Number)))
		PC++
		goto start

	case Div:
		push(pop().(a.Number).Div(pop().(a.Number)))
		PC++
		goto start

	case Mod:
		push(pop().(a.Number).Mod(pop().(a.Number)))
		PC++
		goto start

	case Eq:
		if pop().(a.Number).Cmp(pop().(a.Number)) == a.EqualTo {
			push(a.True)
			PC++
			goto start
		}
		push(a.False)
		PC++
		goto start

	case Neq:
		if pop().(a.Number).Cmp(pop().(a.Number)) != a.EqualTo {
			push(a.True)
			PC++
			goto start
		}
		push(a.False)
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
		PC = INST[PC].Op1
		goto start

	case Jump:
		PC = INST[PC].Op1
		goto start

	case Return:
		return pop()

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

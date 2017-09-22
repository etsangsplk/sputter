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

	for {
		switch INST[PC].OpCode {
		case NoOp:
			// Do Nothing

		case Pop:
			pop()

		case Load:
			push(LOCALS[INST[PC].Op1])

		case Store:
			LOCALS[INST[PC].Op1] = pop()

		case StoreConst:
			LOCALS[INST[PC].Op2] = DATA[INST[PC].Op1]

		case Clear:
			LOCALS[INST[PC].Op1] = a.Nil

		case Dup:
			STACK[SP] = STACK[SP+1]
			SP--

		case Swap:
			u1 := SP + 1
			u2 := SP + 2
			r1 := STACK[u1]
			STACK[u1] = STACK[u2]
			STACK[u2] = r1

		case Nil:
			push(a.Nil)

		case EmptyList:
			push(a.EmptyList)

		case True:
			push(a.True)

		case False:
			push(a.False)

		case Zero:
			push(a.Zero)

		case One:
			push(a.One)

		case NegOne:
			push(negOne)

		case Const:
			push(DATA[INST[PC].Op1])

		case Def:
			r1 := pop()
			a.GetContextNamespace(c).Put(a.AssertUnqualified(pop()).Name(), r1)

		case Let:
			c1 := LOCALS[INST[PC].Op1].(a.Context)
			r1 := pop()
			c1.Put(a.AssertUnqualified(pop()).Name(), r1)

		case Eval:
			r1, _ := a.MacroExpand(c, pop())
			if e1, b1 := r1.(a.Evaluable); b1 {
				push(e1.Eval(c))
			} else {
				push(r1)
			}

		case Apply:
			s1 := pop().(a.Sequence)
			f1 := pop().(a.Applicable)
			push(f1.Apply(c, s1))

		case Call:
			u1 := INST[PC].Op1
			u2 := SP + u1 + 1
			v1 := make(a.Values, u1)
			copy(v1, STACK[SP+1:u2])
			STACK[u2] = STACK[u2].(a.Applicable).Apply(c, v1)
			SP = u2 - 1

		case Vector:
			u1 := INST[PC].Op1
			u2 := SP + u1
			v1 := make(a.Values, u1)
			copy(v1, STACK[SP+1:u2+1])
			STACK[u2] = v1
			SP = u2 - 1

		case IsSeq:
			if s1, b1 := pop().(a.Sequence); b1 && s1.IsSequence() {
				push(a.True)
			} else {
				push(a.False)
			}

		case First:
			push(pop().(a.Sequence).First())

		case Rest:
			push(pop().(a.Sequence).Rest())

		case Split:
			var r1 a.Value
			if s1, b1 := pop().(a.Sequence); b1 {
				if r1, s1, b1 = s1.Split(); b1 {
					push(s1)
					push(r1)
					push(a.True)
					PC++
					continue
				}
			}
			push(a.False)

		case Prepend:
			r1 := pop()
			push(pop().(a.Sequence).Prepend(r1))

		case Inc:
			u1 := SP + 1
			STACK[u1] = STACK[u1].(a.Number).Add(a.One)

		case Dec:
			u1 := SP + 1
			STACK[u1] = STACK[u1].(a.Number).Sub(a.One)

		case Add:
			push(pop().(a.Number).Add(pop().(a.Number)))

		case Sub:
			push(pop().(a.Number).Sub(pop().(a.Number)))

		case Mul:
			push(pop().(a.Number).Mul(pop().(a.Number)))

		case Div:
			push(pop().(a.Number).Div(pop().(a.Number)))

		case Mod:
			push(pop().(a.Number).Mod(pop().(a.Number)))

		case Eq:
			if pop().(a.Number).Cmp(pop().(a.Number)) == a.EqualTo {
				push(a.True)
			} else {
				push(a.False)
			}

		case Neq:
			if pop().(a.Number).Cmp(pop().(a.Number)) != a.EqualTo {
				push(a.True)
			} else {
				push(a.False)
			}

		case Gt:
			if pop().(a.Number).Cmp(pop().(a.Number)) == a.GreaterThan {
				push(a.True)
			} else {
				push(a.False)
			}

		case Gte:
			cmp := pop().(a.Number).Cmp(pop().(a.Number))
			if cmp == a.GreaterThan || cmp == a.EqualTo {
				push(a.True)
			} else {
				push(a.False)
			}

		case Lt:
			if pop().(a.Number).Cmp(pop().(a.Number)) == a.LessThan {
				push(a.True)
			} else {
				push(a.False)
			}

		case Lte:
			cmp := pop().(a.Number).Cmp(pop().(a.Number))
			if cmp == a.LessThan || cmp == a.EqualTo {
				push(a.True)
			} else {
				push(a.False)
			}

		case CondJump:
			r1 := pop()
			if r1 != a.False && r1 != a.Nil {
				PC = INST[PC].Op1
				continue
			}

		case Jump:
			PC = INST[PC].Op1
			continue

		case Return:
			return pop()

		case Panic:
			panic(pop())

		default:
			panic("how did we get here?")
		}
		PC++
	}
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

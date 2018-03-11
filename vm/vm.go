package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

type (
	// Module is the basic translation unit for the VM
	Module struct {
		a.BaseFunction
		LocalsSize   uint
		Data         a.Values
		Instructions Instructions
	}
)

var negOne = a.Zero.Sub(a.One)

// Apply makes Module applicable
func (m *Module) Apply(c a.Context, args a.Sequence) a.Value {
	// Registers
	var PC uint
	DATA := m.Data
	INST := m.Instructions

	// Note: Heap allocated locals slow (fib-vm 30) down by at least 20%
	LOCALS := make(a.Values, m.LocalsSize)
	LOCALS[Context] = c
	LOCALS[Args] = args

	for {
		inst := &INST[PC]
		switch inst.OpCode {
		case Const:
			LOCALS[inst.Op2] = DATA[inst.Op1]

		case Dup:
			LOCALS[inst.Op2] = LOCALS[inst.Op1]

		case Nil:
			LOCALS[inst.Op1] = a.Nil

		case EmptyList:
			LOCALS[inst.Op1] = a.EmptyList

		case True:
			LOCALS[inst.Op1] = a.True

		case False:
			LOCALS[inst.Op1] = a.False

		case Zero:
			LOCALS[inst.Op1] = a.Zero

		case One:
			LOCALS[inst.Op1] = a.One

		case NegOne:
			LOCALS[inst.Op1] = negOne

		case NamespacePut:
			LOCALS[inst.Op3].(a.Namespace).Put(
				LOCALS[inst.Op2].(a.LocalSymbol).Name(),
				LOCALS[inst.Op1])

		case Let:
			LOCALS[inst.Op3].(a.Context).Put(
				LOCALS[inst.Op2].(a.LocalSymbol).Name(),
				LOCALS[inst.Op1])

		case Eval:
			v1, _ := a.MacroExpand(c, LOCALS[inst.Op1])
			if e1, b1 := v1.(a.Evaluable); b1 {
				LOCALS[inst.Op2] = e1.Eval(c)
			} else {
				LOCALS[inst.Op2] = v1
			}

		case Apply:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op2].(a.Applicable).
					Apply(c, LOCALS[inst.Op1].(a.Sequence))

		case Vector:
			s1 := inst.Op1
			e1 := inst.Op2
			v1 := make(a.Values, e1-s1)
			copy(v1, LOCALS[s1:e1])
			LOCALS[inst.Op3] = v1

		case IsSeq:
			v1 := LOCALS[inst.Op1]
			if s1, b1 := v1.(a.Sequence); b1 && s1.IsSequence() {
				LOCALS[inst.Op2] = a.True
			} else {
				LOCALS[inst.Op2] = a.False
			}

		case First:
			LOCALS[inst.Op2] = LOCALS[inst.Op1].(a.Sequence).First()

		case Rest:
			LOCALS[inst.Op2] = LOCALS[inst.Op1].(a.Sequence).Rest()

		case Split:
			if s1, b1 := LOCALS[inst.Op1].(a.Sequence); b1 {
				if f1, r1, b2 := s1.Split(); b2 {
					LOCALS[inst.Op2] = a.True
					LOCALS[inst.Op3] = f1
					LOCALS[inst.Op4] = r1
					break
				}
			}
			LOCALS[inst.Op2] = a.False

		case Prepend:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op2].(a.Sequence).Prepend(LOCALS[inst.Op1])

		case Inc:
			LOCALS[inst.Op2] = LOCALS[inst.Op1].(a.Number).Add(a.One)

		case Dec:
			LOCALS[inst.Op2] = LOCALS[inst.Op1].(a.Number).Sub(a.One)

		case Add:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op1].(a.Number).
					Add(LOCALS[inst.Op2].(a.Number))

		case Sub:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op1].(a.Number).
					Sub(LOCALS[inst.Op2].(a.Number))

		case Mul:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op1].(a.Number).
					Mul(LOCALS[inst.Op2].(a.Number))

		case Div:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op1].(a.Number).
					Div(LOCALS[inst.Op2].(a.Number))

		case Mod:
			LOCALS[inst.Op3] =
				LOCALS[inst.Op1].(a.Number).
					Mod(LOCALS[inst.Op2].(a.Number))

		case Eq:
			n1 := LOCALS[inst.Op1].(a.Number)
			n2 := LOCALS[inst.Op2].(a.Number)
			if n1.Cmp(n2) == a.EqualTo {
				LOCALS[inst.Op3] = a.True
			} else {
				LOCALS[inst.Op3] = a.False
			}

		case Neq:
			n1 := LOCALS[inst.Op1].(a.Number)
			n2 := LOCALS[inst.Op2].(a.Number)
			if n1.Cmp(n2) != a.EqualTo {
				LOCALS[inst.Op3] = a.True
			} else {
				LOCALS[inst.Op3] = a.False
			}

		case Gt:
			n1 := LOCALS[inst.Op1].(a.Number)
			n2 := LOCALS[inst.Op2].(a.Number)
			if n1.Cmp(n2) == a.GreaterThan {
				LOCALS[inst.Op3] = a.True
			} else {
				LOCALS[inst.Op3] = a.False
			}

		case Gte:
			n1 := LOCALS[inst.Op1].(a.Number)
			n2 := LOCALS[inst.Op2].(a.Number)
			cmp := n1.Cmp(n2)
			if cmp == a.GreaterThan || cmp == a.EqualTo {
				LOCALS[inst.Op3] = a.True
			} else {
				LOCALS[inst.Op3] = a.False
			}

		case Lt:
			n1 := LOCALS[inst.Op1].(a.Number)
			n2 := LOCALS[inst.Op2].(a.Number)
			if n1.Cmp(n2) == a.LessThan {
				LOCALS[inst.Op3] = a.True
			} else {
				LOCALS[inst.Op3] = a.False
			}

		case Lte:
			n1 := LOCALS[inst.Op1].(a.Number)
			n2 := LOCALS[inst.Op2].(a.Number)
			cmp := n1.Cmp(n2)
			if cmp == a.LessThan || cmp == a.EqualTo {
				LOCALS[inst.Op3] = a.True
			} else {
				LOCALS[inst.Op3] = a.False
			}

		case CondJump:
			v1 := LOCALS[inst.Op2]
			if v1 != a.False && v1 != a.Nil {
				PC = inst.Op1
				continue
			}

		case Jump:
			PC = inst.Op1
			continue

		case Return:
			return LOCALS[inst.Op1]

		case Panic:
			panic(LOCALS[inst.Op1])

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
		Data:         m.Data,
		Instructions: m.Instructions,
	}
}

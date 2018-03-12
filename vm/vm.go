package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

type (
	// Module is the basic translation unit for the VM
	Module struct {
		a.BaseFunction
		VarCount     uint
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

	// Note: Heap allocated vars slow (fib-vm 30) down by at least 20%
	VARS := make(a.Values, m.VarCount)
	VARS[Context] = c
	VARS[Args] = args

	for {
		i := &INST[PC]
		switch i.OpCode {
		case Const:
			VARS[i.Op2] = DATA[i.Op1]

		case Dup:
			VARS[i.Op2] = VARS[i.Op1]

		case Nil:
			VARS[i.Op1] = a.Nil

		case EmptyList:
			VARS[i.Op1] = a.EmptyList

		case True:
			VARS[i.Op1] = a.True

		case False:
			VARS[i.Op1] = a.False

		case Zero:
			VARS[i.Op1] = a.Zero

		case One:
			VARS[i.Op1] = a.One

		case NegOne:
			VARS[i.Op1] = negOne

		case NamespacePut:
			VARS[i.Op3].(a.Namespace).Put(
				VARS[i.Op2].(a.LocalSymbol).Name(), VARS[i.Op1])

		case Let:
			VARS[i.Op3].(a.Context).Put(
				VARS[i.Op2].(a.LocalSymbol).Name(), VARS[i.Op1])

		case Eval:
			v1, _ := a.MacroExpand(c, VARS[i.Op1])
			if e1, b1 := v1.(a.Evaluable); b1 {
				VARS[i.Op2] = e1.Eval(c)
			} else {
				VARS[i.Op2] = v1
			}

		case Apply:
			VARS[i.Op3] = VARS[i.Op2].(a.Applicable).
				Apply(c, VARS[i.Op1].(a.Sequence))

		case Vector:
			s1 := i.Op1
			e1 := i.Op2
			v1 := make(a.Values, e1-s1)
			copy(v1, VARS[s1:e1])
			VARS[i.Op3] = v1

		case IsSeq:
			v1 := VARS[i.Op1]
			if s1, b1 := v1.(a.Sequence); b1 && s1.IsSequence() {
				VARS[i.Op2] = a.True
			} else {
				VARS[i.Op2] = a.False
			}

		case First:
			VARS[i.Op2] = VARS[i.Op1].(a.Sequence).First()

		case Rest:
			VARS[i.Op2] = VARS[i.Op1].(a.Sequence).Rest()

		case Split:
			if s1, b1 := VARS[i.Op1].(a.Sequence); b1 {
				if f1, r1, b2 := s1.Split(); b2 {
					VARS[i.Op2] = a.True
					VARS[i.Op3] = f1
					VARS[i.Op4] = r1
					break
				}
			}
			VARS[i.Op2] = a.False

		case Prepend:
			VARS[i.Op3] = VARS[i.Op2].(a.Sequence).Prepend(VARS[i.Op1])

		case Inc:
			VARS[i.Op2] = VARS[i.Op1].(a.Number).Add(a.One)

		case Dec:
			VARS[i.Op2] = VARS[i.Op1].(a.Number).Sub(a.One)

		case Add:
			VARS[i.Op3] = VARS[i.Op1].(a.Number).Add(VARS[i.Op2].(a.Number))

		case Sub:
			VARS[i.Op3] = VARS[i.Op1].(a.Number).Sub(VARS[i.Op2].(a.Number))

		case Mul:
			VARS[i.Op3] = VARS[i.Op1].(a.Number).Mul(VARS[i.Op2].(a.Number))

		case Div:
			VARS[i.Op3] = VARS[i.Op1].(a.Number).Div(VARS[i.Op2].(a.Number))

		case Mod:
			VARS[i.Op3] = VARS[i.Op1].(a.Number).Mod(VARS[i.Op2].(a.Number))

		case Eq:
			n1 := VARS[i.Op1].(a.Number)
			n2 := VARS[i.Op2].(a.Number)
			if n1.Cmp(n2) == a.EqualTo {
				VARS[i.Op3] = a.True
			} else {
				VARS[i.Op3] = a.False
			}

		case Neq:
			n1 := VARS[i.Op1].(a.Number)
			n2 := VARS[i.Op2].(a.Number)
			if n1.Cmp(n2) != a.EqualTo {
				VARS[i.Op3] = a.True
			} else {
				VARS[i.Op3] = a.False
			}

		case Gt:
			n1 := VARS[i.Op1].(a.Number)
			n2 := VARS[i.Op2].(a.Number)
			if n1.Cmp(n2) == a.GreaterThan {
				VARS[i.Op3] = a.True
			} else {
				VARS[i.Op3] = a.False
			}

		case Gte:
			n1 := VARS[i.Op1].(a.Number)
			n2 := VARS[i.Op2].(a.Number)
			cmp := n1.Cmp(n2)
			if cmp == a.GreaterThan || cmp == a.EqualTo {
				VARS[i.Op3] = a.True
			} else {
				VARS[i.Op3] = a.False
			}

		case Lt:
			n1 := VARS[i.Op1].(a.Number)
			n2 := VARS[i.Op2].(a.Number)
			if n1.Cmp(n2) == a.LessThan {
				VARS[i.Op3] = a.True
			} else {
				VARS[i.Op3] = a.False
			}

		case Lte:
			n1 := VARS[i.Op1].(a.Number)
			n2 := VARS[i.Op2].(a.Number)
			cmp := n1.Cmp(n2)
			if cmp == a.LessThan || cmp == a.EqualTo {
				VARS[i.Op3] = a.True
			} else {
				VARS[i.Op3] = a.False
			}

		case CondJump:
			v1 := VARS[i.Op2]
			if v1 != a.False && v1 != a.Nil {
				PC = i.Op1
				continue
			}

		case Jump:
			PC = i.Op1
			continue

		case Return:
			return VARS[i.Op1]

		case Panic:
			panic(VARS[i.Op1])

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
		VarCount:     m.VarCount,
		Data:         m.Data,
		Instructions: m.Instructions,
	}
}

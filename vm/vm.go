package vm

// Make it fast, and then make it fast

import a "github.com/kode4food/sputter/api"

// MaxValueRegisters is the maximum number of Value-specific registers
const MaxValueRegisters = 8

type (
	// Module is the basic translation unit for the VM
	Module struct {
		a.BaseFunction

		Data         a.Values
		Instructions Instructions
	}
)

var negOne = a.Zero.Sub(a.One)

// Apply makes Module applicable
func (m *Module) Apply(c a.Context, args a.Sequence) a.Value {
	// Registers
	var PC uint
	INST := m.Instructions

	var VAL [MaxValueRegisters]a.Value
	VAL[Context] = c
	VAL[Args] = args

	for {
		i := &INST[PC]
		switch i.OpCode {
		case Const:
			VAL[i.Op2] = m.Data[i.Op1]

		case Dup:
			VAL[i.Op2] = VAL[i.Op1]

		case Nil:
			VAL[i.Op1] = a.Nil

		case EmptyList:
			VAL[i.Op1] = a.EmptyList

		case True:
			VAL[i.Op1] = a.True

		case False:
			VAL[i.Op1] = a.False

		case Zero:
			VAL[i.Op1] = a.Zero

		case One:
			VAL[i.Op1] = a.One

		case NegOne:
			VAL[i.Op1] = negOne

		case NamespacePut:
			VAL[i.Op3].(a.Namespace).Put(
				VAL[i.Op2].(a.LocalSymbol).Name(), VAL[i.Op1])

		case Let:
			VAL[i.Op3].(a.Context).Put(
				VAL[i.Op2].(a.LocalSymbol).Name(), VAL[i.Op1])

		case Eval:
			v1, _ := a.MacroExpand(c, VAL[i.Op1])
			if e1, b1 := v1.(a.Evaluable); b1 {
				VAL[i.Op2] = e1.Eval(c)
			} else {
				VAL[i.Op2] = v1
			}

		case Apply:
			VAL[i.Op3] = VAL[i.Op2].(a.Applicable).
				Apply(c, VAL[i.Op1].(a.Sequence))

		case Vector:
			s1 := i.Op1
			e1 := i.Op2
			v1 := make(a.Values, e1-s1)
			copy(v1, VAL[s1:e1])
			VAL[i.Op3] = v1

		case IsSeq:
			v1 := VAL[i.Op1]
			if s1, b1 := v1.(a.Sequence); b1 && s1.IsSequence() {
				VAL[i.Op2] = a.True
			} else {
				VAL[i.Op2] = a.False
			}

		case First:
			VAL[i.Op2] = VAL[i.Op1].(a.Sequence).First()

		case Rest:
			VAL[i.Op2] = VAL[i.Op1].(a.Sequence).Rest()

		case Split:
			if s1, b1 := VAL[i.Op1].(a.Sequence); b1 {
				if f1, r1, b2 := s1.Split(); b2 {
					VAL[i.Op2] = a.True
					VAL[i.Op3] = f1
					VAL[i.Op4] = r1
					break
				}
			}
			VAL[i.Op2] = a.False

		case Prepend:
			VAL[i.Op3] = VAL[i.Op2].(a.Sequence).Prepend(VAL[i.Op1])

		case Inc:
			VAL[i.Op2] = VAL[i.Op1].(a.Number).Add(a.One)

		case Dec:
			VAL[i.Op2] = VAL[i.Op1].(a.Number).Sub(a.One)

		case Add:
			VAL[i.Op3] = VAL[i.Op1].(a.Number).Add(VAL[i.Op2].(a.Number))

		case Sub:
			VAL[i.Op3] = VAL[i.Op1].(a.Number).Sub(VAL[i.Op2].(a.Number))

		case Mul:
			VAL[i.Op3] = VAL[i.Op1].(a.Number).Mul(VAL[i.Op2].(a.Number))

		case Div:
			VAL[i.Op3] = VAL[i.Op1].(a.Number).Div(VAL[i.Op2].(a.Number))

		case Mod:
			VAL[i.Op3] = VAL[i.Op1].(a.Number).Mod(VAL[i.Op2].(a.Number))

		case Eq:
			n1 := VAL[i.Op1].(a.Number)
			n2 := VAL[i.Op2].(a.Number)
			if n1.Cmp(n2) == a.EqualTo {
				VAL[i.Op3] = a.True
			} else {
				VAL[i.Op3] = a.False
			}

		case Neq:
			n1 := VAL[i.Op1].(a.Number)
			n2 := VAL[i.Op2].(a.Number)
			if n1.Cmp(n2) != a.EqualTo {
				VAL[i.Op3] = a.True
			} else {
				VAL[i.Op3] = a.False
			}

		case Gt:
			n1 := VAL[i.Op1].(a.Number)
			n2 := VAL[i.Op2].(a.Number)
			if n1.Cmp(n2) == a.GreaterThan {
				VAL[i.Op3] = a.True
			} else {
				VAL[i.Op3] = a.False
			}

		case Gte:
			n1 := VAL[i.Op1].(a.Number)
			n2 := VAL[i.Op2].(a.Number)
			cmp := n1.Cmp(n2)
			if cmp == a.GreaterThan || cmp == a.EqualTo {
				VAL[i.Op3] = a.True
			} else {
				VAL[i.Op3] = a.False
			}

		case Lt:
			n1 := VAL[i.Op1].(a.Number)
			n2 := VAL[i.Op2].(a.Number)
			if n1.Cmp(n2) == a.LessThan {
				VAL[i.Op3] = a.True
			} else {
				VAL[i.Op3] = a.False
			}

		case Lte:
			n1 := VAL[i.Op1].(a.Number)
			n2 := VAL[i.Op2].(a.Number)
			cmp := n1.Cmp(n2)
			if cmp == a.LessThan || cmp == a.EqualTo {
				VAL[i.Op3] = a.True
			} else {
				VAL[i.Op3] = a.False
			}

		case CondJump:
			v1 := VAL[i.Op2]
			if v1 != a.False && v1 != a.Nil {
				PC = i.Op1
				continue
			}

		case Jump:
			PC = i.Op1
			continue

		case Return:
			return VAL[i.Op1]

		case Panic:
			panic(VAL[i.Op1])

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
		Data:         m.Data,
		Instructions: m.Instructions,
	}
}

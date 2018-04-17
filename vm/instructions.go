package vm

// Instruction Argument Types
const (
	Read InstructionArgType = iota
	Write
	Data
	Index
	Offset
	LabelID
)

type (
	// Instruction represents a decoded VM instruction
	Instruction struct {
		OpCode OpCode
		Op1    uint
		Op2    uint
		Op3    uint
		Op4    uint
	}

	// Instructions represents multiple VM instructions
	Instructions []Instruction

	// InstructionArgType represents the type of an instruction argument
	InstructionArgType int

	// InstructionType holds information about a valid Instruction
	InstructionType struct {
		OpCode   OpCode
		ArgCount uint
		ArgTypes []InstructionArgType
	}
)

// InstructionTypes describes the Instructions we understand
var InstructionTypes = []InstructionType{
	MakeInstType(Label, LabelID),
	MakeInstType(Const, Data, Write),
	MakeInstType(Dup, Read, Write),
	MakeInstType(Nil, Write),
	MakeInstType(EmptyList, Write),
	MakeInstType(True, Write),
	MakeInstType(False, Write),
	MakeInstType(Zero, Write),
	MakeInstType(One, Write),
	MakeInstType(NegOne, Write),
	MakeInstType(NamespacePut, Read, Read, Write),
	MakeInstType(Let, Read, Read, Write),
	MakeInstType(Eval, Read, Write),
	MakeInstType(Apply, Read, Read, Write),
	MakeInstType(Vector, Read, Read, Write),
	MakeInstType(GetValue, Read, Read, Write),
	MakeInstType(GetArg, Index, Write),
	MakeInstType(IsSeq, Read, Write),
	MakeInstType(First, Read, Write),
	MakeInstType(Rest, Read, Write),
	MakeInstType(Split, Read, Write, Write, Write),
	MakeInstType(Prepend, Read, Read, Write),
	MakeInstType(Inc, Read, Write),
	MakeInstType(Dec, Read, Write),
	MakeInstType(Add, Read, Read, Write),
	MakeInstType(Sub, Read, Read, Write),
	MakeInstType(Mul, Read, Read, Write),
	MakeInstType(Div, Read, Read, Write),
	MakeInstType(Mod, Read, Read, Write),
	MakeInstType(Eq, Read, Read, Write),
	MakeInstType(Neq, Read, Read, Write),
	MakeInstType(Gt, Read, Read, Write),
	MakeInstType(Gte, Read, Read, Write),
	MakeInstType(Lt, Read, Read, Write),
	MakeInstType(Lte, Read, Read, Write),
	MakeInstType(CondJumpLabel, Read, LabelID),
	MakeInstType(CondJump, Read, Offset),
	MakeInstType(JumpLabel, LabelID),
	MakeInstType(Jump, Offset),
	MakeInstType(Return, Read),
	MakeInstType(Panic, Read),
}

// MakeInstType constructs an Instruction Type descriptor in a cleaner way
func MakeInstType(o OpCode, a ...InstructionArgType) InstructionType {
	return InstructionType{
		OpCode:   o,
		ArgCount: uint(len(a)),
		ArgTypes: a,
	}
}

// MakeInst constructs Instructions in a cleaner way
func MakeInst(o OpCode, args ...uint) Instruction {
	i := Instruction{
		OpCode: o,
	}
	switch len(args) {
	case 4:
		i.Op4 = args[3]
		fallthrough
	case 3:
		i.Op3 = args[2]
		fallthrough
	case 2:
		i.Op2 = args[1]
		fallthrough
	case 1:
		i.Op1 = args[0]
		fallthrough
	case 0:
		return i
	default:
		panic("error: call to MakeInst has too many operands")
	}
}

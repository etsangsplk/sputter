package vm

// OpCode represents a decoded VM instruction identifier
type OpCode uint

// These are the OpCodes recognized by the Virtual Machine
const (
	Const OpCode = iota
	Dup
	Nil
	EmptyList
	True
	False
	Zero
	One
	NegOne
	NamespacePut
	Let
	Eval
	Apply
	Vector
	IsSeq
	First
	Rest
	Split
	Prepend
	Inc
	Dec
	Add
	Sub
	Mul
	Div
	Mod
	Eq
	Neq
	Gt
	Gte
	Lt
	Lte
	CondJump
	Jump
	Return
	Panic
	// Should never reach the VM
	Label
	CondJumpLabel
	JumpLabel
)

// These are the default positions for certain registers
const (
	Context uint = iota
	Args
	Vars
)
